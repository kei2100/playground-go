package main

import (
	"fmt"
)

type Message string

type Greeter struct {
	Message Message
}

type Event struct {
	Greeter Greeter
}

func NewMessage() Message {
	return Message("hello wire")
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (g Greeter) Greet() Message {
	return g.Message
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	e := InitializeEvent()
	e.Start()
}
