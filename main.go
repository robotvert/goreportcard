package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/robotvert/goreportcard/handlers"
)

func main() {
	resp, err := handlers.CheckPackage(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	pp.Println(resp)
}
