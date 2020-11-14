package main

import "syscall/js"

func main() {
	doc := js.Global().Get("document")
	before := doc.Call("getElementById", "before")
	before.Set("value", "{}")
}
