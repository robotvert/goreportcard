# Go Report Card - library style

The upstream project is a full fledged web app but I only need to run it locally, on a given package's source code.

This Quick n' Dirty™ variant keeps only the bare minimum code (no mongo caching, no http handlers etc.).
Some changes were made accordingly, while trying to keep those minimal.


## SETUP

To run properly several tools need to be installed, using the following commands:

```
go get -u github.com/fzipp/gocyclo

go get -u golang.org/x/tools/cmd/vet

go get -u github.com/golang/lint/golint
```

## SAMPLE USAGE

```
package main

import (
  "fmt"
  "os"
  "github.com/robotvert/goreportcard"
)

func main() {
  resp, err := goreportcard.CheckPackage(os.Args[1])
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(resp.Grade)
}
```
