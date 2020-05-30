package main

import (
	"fmt"
)

func getWASM() {
	registerCallbacks()
}

func main() {
	c := make(chan struct{}, 0)

	getWASM()
	fmt.Println("Hello wasm!!")
	<-c
}
