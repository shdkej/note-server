package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"syscall/js"
)

func initial() {
	document := js.Global().Get("document")
	p := document.Call("getElementById", "article")
	path := "recommend.txt"
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	u.Path = path
	response, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	p.Set("innerText", string(data))
}

/*
func add(this js.Value, i []js.Value) interface{} {
	js.Global().Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
	result := js.ValueOf(i[0].Int() + i[1].Int()).String()
	println(result)
	return result
}
*/

func handleRect(this js.Value, i []js.Value) interface{} {
	document := js.Global().Get("document")
	p := document.Call("getElementById", "square")
	p.Set("innerText", "notest")
	return p
}

func toggleInput(this js.Value, i []js.Value) interface{} {
	document := js.Global().Get("document")
	p := document.Call("getElementById", "input1")
	p.Set("style", "display:block;")
	return p
}

func toggleNavbar(this js.Value, i []js.Value) interface{} {
	document := js.Global().Get("document")
	p := document.Call("getElementById", "go_wasm").
		Get("classList").
		Call("toggle", "active")
	return p
}

func storeInput(this js.Value, i []js.Value) interface{} {
	document := js.Global().Get("document")
	p := document.Call("getElementById", "article")

	go func() {
		//path := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String() + ".md"
		path := "spring.md"
		href := js.Global().Get("location").Get("href")
		u, err := url.Parse(href.String())
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		u.Path = path
		response, err := http.Get(u.String())
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		p.Set("innerText", string(data))
	}()
	return p
}

func registerCallbacks() {
	js.Global().Set("handleRect", js.FuncOf(handleRect))
	js.Global().Set("toggleInput", js.FuncOf(toggleInput))
	js.Global().Set("storeInput", js.FuncOf(storeInput))
	js.Global().Set("toggleNavbar", js.FuncOf(toggleNavbar))
}
