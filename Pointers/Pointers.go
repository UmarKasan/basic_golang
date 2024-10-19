package main

import "fmt"

// This project is called a linked list using pointers
type greeting struct {
	text     string
	next     *greeting
	previous *greeting
}

func sayGreeting(say *greeting) {
	if say == nil { // Need to add as Println will generate error if nil is not added
		return
	}
	fmt.Println(say.text)
	sayGreeting(say.next)
	sayGreeting(say.previous)
	fmt.Println(say.text)
}

func main() {
	greeting1 := &greeting{text: "Hello"}
	greeting2 := &greeting{text: "konnichiwa"}
	greeting3 := &greeting{text: "Bonjour"}

	greeting1.previous = greeting3
	greeting1.next = greeting2
	greeting2.previous = greeting1
	greeting2.next = greeting3
	greeting3.previous = greeting2

	sayGreeting(greeting1)
}
