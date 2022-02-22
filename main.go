package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data,err:=ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
