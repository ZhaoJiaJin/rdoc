# ztcd
a k-v document database based on raft

# features

* Distributed lock(with timeout)
    * Lock 
    * Unlock
* Read
    * FindAll
    * SpecifyConditions
    
* Create
* Update
* Delete

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

