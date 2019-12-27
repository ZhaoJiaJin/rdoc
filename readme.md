# ztcd
a k-v document database based on raft

# features

* document have auto generated id
* Distributed lock(with timeout)
    * Lock 
    * Unlock
* Read
    * FindAll
    * SpecifyConditions
    
* Create
    * Add a document
* Update
    * SpecifyConditions to update documents
* Delete
    * DeleteAll
    * SpecifyConditions to delete documents

# how to SpecifyConditions
    * Equality
    ```python
    {key:value}
    ```
    * Less Than
    ```python
    {
        key:{"$lt":value}
    }
    ```
    * Great Than
    ```python
    {
        key:{"$gt":value}
    }
    ```
    * Less Than Equals
    ```python
    {
        key:{"$lte":value}
    }
    ```
    * Great Than Equals
    ```python
    {
        key:{"$gte":value}
    }
    ```
    * Not Equals
    ```python
    {
        key:{"$ne":value}
    }
    ```

