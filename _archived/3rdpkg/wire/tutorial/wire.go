//go:generate wire
//+build wireinject

package main

import "github.com/google/wire"

// InitializeEvent initializes a Event.
// In Wire parlance, InitializeEvent is an "injector."
func InitializeEvent() Event {
	// Rather than go through the trouble of initializing each component in turn and passing it into the next one,
	// we instead have a single call to wire.Build passing in the initializers we want to use.
	// In Wire, initializers are known as "providers," functions which provide a particular type.
	//wire.Build(NewEvent, NewGreeter, NewMessage)
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}
}
