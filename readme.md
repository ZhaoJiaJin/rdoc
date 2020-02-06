# A high available document database based on Raft Algorithm, inspired by [tiedot](https://github.com/HouzuoGuo/tiedot/wiki/Tutorial) and [Raft](https://raft.github.io/)


# how to design a document database.

* Doc:
* Col: one collection contains many documents
* index:
  * index
  * id to phyid



# basic data structure

* Collection: every collection should be a separated array, this array will contains multiple documents. the initial size of array should be SizeOfASingleDocument*#MaxDocument. If







# Question

1. do we need to support add en entry to an array in a doc.

