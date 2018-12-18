package main

import (
	"fmt"
)

// Message type
type Message string

// Greeter struct
type Greeter struct {
	Message Message
}

// Event strut
type Event struct {
	Greeter Greeter
}

// NewMessage provides a new message
func NewMessage() Message {
	return Message("hello wire")
}

// NewGreeter provides a new Greeter
func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

// NewEvent provides a new Event
func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

// Greet Message
func (g Greeter) Greet() Message {
	return g.Message
}

// Start event
func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	e := InitializeEvent()
	e.Start()
}
