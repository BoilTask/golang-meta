package event

import (
	"log/slog"
	"reflect"
	"testing"
)

type TestEvent struct {
	Event
}

func TestBase(t *testing.T) {
	l := NewListenerDefault(
		func(event reflect.Type, p ...Payload) {
			slog.Info("called test event", "event", event, "p", p)
		},
	)
	RegisterListener[TestEvent](l)
	Invoke[TestEvent]("test")
	UnregisterListener[TestEvent](l)
	Invoke[TestEvent]()
}
