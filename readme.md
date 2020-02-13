# A high available document database based on Raft Algorithm, inspired by [tiedot](https://github.com/HouzuoGuo/tiedot/wiki/Tutorial) and [Raft](https://raft.github.io/)


# API

## collection

### create collection

```shell

curl --location --request GET 'localhost:8080/create?col=Feeds'

```

```golang

package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "localhost:8080/create?col=Feeds"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
  }
  res, err := client.Do(req)
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)

  fmt.Println(string(body))
}

```
