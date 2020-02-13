// Query processor.

package db

import (
	"errors"
	"fmt"
	"strings"
    "rdoc/utils"
)

// EvalUnion union operation
func EvalUnion(exprs []interface{}, src *Col, result *map[string]struct{}) (err error) {
	for _, subExpr := range exprs {
		if err = evalQuery(subExpr, src, result, false); err != nil {
			return
		}
	}
	return
}

// EvalAllIDs get all document id
func EvalAllIDs(src *Col, result *map[string]struct{}) (err error) {
	src.ForEachDoc(func(id string, _ *Doc) bool {
		(*result)[id] = struct{}{}
		return true
	})
	return
}

// Value equity check ("attribute == value") using hash lookup.
func Lookup(lookupValue interface{}, expr map[string]interface{}, src *Col, result *map[string]struct{}) (err error) {
	// Figure out lookup path - JSON array "in"
	path, hasPath := expr["in"]
	if !hasPath {
		return errors.New("Missing lookup path `in`")
	}
	vecPath := make([]string, 0)
	if vecPathInterface, ok := path.([]interface{}); ok {
		for _, v := range vecPathInterface {
			vecPath = append(vecPath, fmt.Sprint(v))
		}
	} else {
		return fmt.Errorf("Expecting vector lookup path `in`, but %v given", path)
	}
	// Figure out result number limit
	intLimit := int(0)
	if limit, hasLimit := expr["limit"]; hasLimit {
		if floatLimit, ok := limit.(float64); ok {
			intLimit = int(floatLimit)
		} else if _, ok := limit.(int); ok {
			intLimit = limit.(int)
		} else {
			return ErrLimit
		}
	}
	lookupStrValue := fmt.Sprint(lookupValue) // the value to look for
	lookupValueHash := utils.StrHash(lookupStrValue)
	scanPath := strings.Join(vecPath, INDEX_PATH_SEP)
	/*if _, indexed := src.indexPaths[scanPath]; !indexed {
		// can not use eq unless indexed
		return dberr.New(ErrNotIDX, scanPath, expr)
	}*/
	//num := lookupValueHash % src.db.numParts
	//ht := src.hts[num][scanPath]
	//ht.Lock.RLock()
	vals,err := src.Query(scanPath, lookupValueHash, intLimit)
    if err != nil{
        return err
    }
	//ht.Lock.RUnlock()
	for _, match := range vals {
		// Filter result to avoid hash collision
		if doc := src.ReadDoc(match); doc != nil {
			for _, v := range GetIn(doc.data, vecPath) {
				if fmt.Sprint(v) == lookupStrValue {
					(*result)[match] = struct{}{}
				}
			}
		}
	}
	return
}

// PathExistence existence check (value != nil) 
func PathExistence(hasPath interface{}, expr map[string]interface{}, src *Col, result *map[string]struct{}) (err error) {
	// Figure out the path
	vecPath := make([]string, 0)
	if vecPathInterface, ok := hasPath.([]interface{}); ok {
		for _, v := range vecPathInterface {
			vecPath = append(vecPath, fmt.Sprint(v))
		}
	} else {
		return errors.New(fmt.Sprintf("Expecting vector path, but %v given", hasPath))
	}
	// Figure out result number limit
	intLimit := 0
	if limit, hasLimit := expr["limit"]; hasLimit {
		if floatLimit, ok := limit.(float64); ok {
			intLimit = int(floatLimit)
		} else if _, ok := limit.(int); ok {
			intLimit = limit.(int)
		} else {
			return ErrLimitType
		}
	}
	jointPath := strings.Join(vecPath, INDEX_PATH_SEP)
    vals,err := src.QueryExist(jointPath, intLimit)
    if err != nil{
        return err
    }
	//ht.Lock.RUnlock()
	for _, match := range vals {
		// Filter result to avoid hash collision
		if doc := src.ReadDoc(match); doc != nil {
					(*result)[match] = struct{}{}
		}
	}
	return

	return nil
}

// Calculate intersection of sub-query results.
func Intersect(subExprs interface{}, src *Col, result *map[string]struct{}) (err error) {
	myResult := make(map[string]struct{})
	if subExprVecs, ok := subExprs.([]interface{}); ok {
		first := true
		for _, subExpr := range subExprVecs {
			subResult := make(map[string]struct{})
			intersection := make(map[string]struct{})
			if err = evalQuery(subExpr, src, &subResult, false); err != nil {
				return
			}
			if first {
				myResult = subResult
				first = false
			} else {
				for k := range subResult {
					if _, inBoth := myResult[k]; inBoth {
						intersection[k] = struct{}{}
					}
				}
				myResult = intersection
			}
		}
		for docID := range myResult {
			(*result)[docID] = struct{}{}
		}
	} else {
		return ErrNoSubQuery
	}
	return
}

// Calculate complement of sub-query results.
func Complement(subExprs interface{}, src *Col, result *map[string]struct{}) (err error) {
	myResult := make(map[string]struct{})
	if subExprVecs, ok := subExprs.([]interface{}); ok {
		for _, subExpr := range subExprVecs {
			subResult := make(map[string]struct{})
			complement := make(map[string]struct{})
			if err = evalQuery(subExpr, src, &subResult, false); err != nil {
				return
			}
			for k := range subResult {
				if _, inBoth := myResult[k]; !inBoth {
					complement[k] = struct{}{}
				}
			}
			for k := range myResult {
				if _, inBoth := subResult[k]; !inBoth {
					complement[k] = struct{}{}
				}
			}
			myResult = complement
		}
		for docID := range myResult {
			(*result)[docID] = struct{}{}
		}
	} else {
		return ErrNoSubQuery
	}
	return
}

// IntRange Look for indexed integer values within the specified integer range.
func IntRange(intFrom interface{}, expr map[string]interface{}, src *Col, result *map[string]struct{}) (err error) {
	path, hasPath := expr["in"]
	if !hasPath {
		return errors.New("Missing path `in`")
	}
	// Figure out the path
	vecPath := make([]string, 0)
	if vecPathInterface, ok := path.([]interface{}); ok {
		for _, v := range vecPathInterface {
			vecPath = append(vecPath, fmt.Sprint(v))
		}
	} else {
		return errors.New(fmt.Sprintf("Expecting vector path `in`, but %v given", path))
	}
	// Figure out result number limit
	intLimit := int(0)
	if limit, hasLimit := expr["limit"]; hasLimit {
		if floatLimit, ok := limit.(float64); ok {
			intLimit = int(floatLimit)
		} else if _, ok := limit.(int); ok {
			intLimit = limit.(int)
		} else {
			return ErrLimitType
		}
	}
	// Figure out the range ("from" value & "to" value)
	from, to := int(0), int(0)
	if floatFrom, ok := intFrom.(float64); ok {
		from = int(floatFrom)
	} else if _, ok := intFrom.(int); ok {
		from = intFrom.(int)
	} else {
		return ErrRangeType
	}
	if intTo, ok := expr["int-to"]; ok {
		if floatTo, ok := intTo.(float64); ok {
			to = int(floatTo)
		} else if _, ok := intTo.(int); ok {
			to = intTo.(int)
		} else {
			return ErrRangeType
		}
	} else if intTo, ok := expr["int to"]; ok {
		if floatTo, ok := intTo.(float64); ok {
			to = int(floatTo)
		} else if _, ok := intTo.(int); ok {
			to = intTo.(int)
		} else {
			return ErrRangeType
		}
	} else {
		return ErrRangeMiss
	}
	/*if to > from && to-from > 1000 || from > to && from-to > 1000 {
		tdlog.CritNoRepeat("Query %v involves index lookup on more than 1000 values, which can be very inefficient", expr)
	}*/
	counter := int(0) // Number of results already collected
	htPath := strings.Join(vecPath, INDEX_PATH_SEP)
	if indexScan := src.IsIndexed(htPath); !indexScan {
		return ErrNotIDX
	}
	if from < to {
		// Forward scan - from low value to high value
		for lookupValue := from; lookupValue <= to; lookupValue++ {
			lookupStrValue := fmt.Sprint(float64(lookupValue))
			hashValue := utils.StrHash(lookupStrValue)
			vals,err := src.Query(htPath, hashValue, int(intLimit))
            if err != nil{
                return err
            }
			for _, docID := range vals {
				if intLimit > 0 && counter == intLimit {
					break
				}
				counter++
				(*result)[docID] = struct{}{}
			}
		}
	} else {
		// Backward scan - from high value to low value
		for lookupValue := from; lookupValue >= to; lookupValue-- {
			lookupStrValue := fmt.Sprint(float64(lookupValue))
			hashValue := utils.StrHash(lookupStrValue)
			vals,err := src.Query(htPath, hashValue, int(intLimit))
            if err != nil{
                return err
            }
			for _, docID := range vals {
				if intLimit > 0 && counter == intLimit {
					break
				}
				counter++
				(*result)[docID] = struct{}{}
			}
		}
	}
	return
}

func evalQuery(q interface{}, src *Col, result *map[string]struct{}, placeSchemaLock bool) (err error) {
	/*if placeSchemaLock {
		src.db.schemaLock.RLock()
		defer src.db.schemaLock.RUnlock()
	}*/
	switch expr := q.(type) {
	case []interface{}: // [sub query 1, sub query 2, etc]
		return EvalUnion(expr, src, result)
	case string:
		if expr == "all" {
			return EvalAllIDs(src, result)
		}
		// Might be single document number
		docID := expr
		(*result)[docID] = struct{}{}
	case map[string]interface{}:
		if lookupValue, lookup := expr["eq"]; lookup { // eq - lookup
			return Lookup(lookupValue, expr, src, result)
		} else if hasPath, exist := expr["has"]; exist { // has - path existence test
			return PathExistence(hasPath, expr, src, result)
		} else if subExprs, intersect := expr["n"]; intersect { // n - intersection
			return Intersect(subExprs, src, result)
		} else if subExprs, complement := expr["c"]; complement { // c - complement
			return Complement(subExprs, src, result)
		} else if intFrom, htRange := expr["int-from"]; htRange { // int-from, int-to - integer range query
			return IntRange(intFrom, expr, src, result)
		} else if intFrom, htRange := expr["int from"]; htRange { // "int from, "int to" - integer range query - same as above, just without dash
			return IntRange(intFrom, expr, src, result)
		} else {
			return errors.New(fmt.Sprintf("Query %v does not contain any operation (lookup/union/etc)", expr))
		}
	}
	return nil
}

// EvalQuery evaluate query 
func EvalQuery(q interface{}, src *Col, result *map[string]struct{}) (err error) {
	return evalQuery(q, src, result, true)
}


