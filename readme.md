# ztcd
a k-v document database based on raft

# features

* distributed lock
* Read
    * FindAll
    * SpecifyConditions
        * Equality
        ```python
        {key:value}
        ```
* Create
* Update
* Delete
