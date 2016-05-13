package main

import (
	"fmt"
)

type event struct {
	DateTime string
	Policy   policy
	Type     string
	Message  string
}

func (e event) String() string {
	return fmt.Sprintf("[%v] %v - %v: %v", e.DateTime, e.Policy.Name, e.Type, e.Message)
}

type EventContext struct {
	observers []chan event
}

func (ec *EventContext) registerEventObserver(c chan event) {
	ec.observers = append(ec.observers, c)
}

// send a copy of the event to every registered observer
func (ec *EventContext) notify(e event) {
	for _, o := range ec.observers {
		go func(o chan event, e event) {
			// TODO: add timeout handler to remove observer
			o <- e
		}(o, e)
	}
}
