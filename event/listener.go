package event

import (
	"meta/object"
	"reflect"
)

type ListenerInterface interface {
	object.Interface
	OnEventInvoked(event reflect.Type, p ...Payload)
}

type ListenerDefault struct {
	ListenerInterface
	callback func(event reflect.Type, p ...Payload)
}

func (l *ListenerDefault) GetName() string {
	return "ListenerDefault"
}

func (l *ListenerDefault) OnEventInvoked(event reflect.Type, p ...Payload) {
	l.callback(event, p)
}

func NewListenerDefault(callback func(event reflect.Type, p ...Payload)) *ListenerDefault {
	return &ListenerDefault{callback: callback}
}
