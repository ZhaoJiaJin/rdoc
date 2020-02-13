# A high available document database based on Raft Algorithm, inspired by [tiedot](https://github.com/HouzuoGuo/tiedot/wiki/Tutorial) and [Raft](https://raft.github.io/)

# how to run 

1. install go
2. go build -o server main.go
3. ./server

# API

## create collection

```shell

curl --location --request GET 'localhost:8080/create?col=Feeds'

```

## get all collections

```shell
curl --location --request GET 'localhost:8080/all'
```

## insert a document into collection

```shell
curl --location --request POST 'localhost:8080/insert' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------383513606155476017747483' \
--form 'col=Feeds' \
--form 'doc={"a":{"b":1},"a1":123123}'
```

## query docs

### query all docs

```shell

curl --location --request POST 'localhost:8080/query' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------759712145504402840994879' \
--form 'col=Feeds' \
--form 'q=["all"]'

```


### query using equal

```shell
curl --location --request POST 'localhost:8080/query' \
--header 'Content-Type: multipart/form-data; boundary=--------------------------715073562991340840763979' \
--form 'col=Feeds' \
--form 'q={"eq": 1, "in": ["a","b"]}'
```


## create index for a collection

```shell
curl --location --request POST 'localhost:8080/index' \
--form 'col=Feeds' \
--form 'path=a,b'
```

