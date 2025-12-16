package event

import (
	"log/slog"
	metaflag "meta/meta-flag"
	metaformat "meta/meta-format"
	metastring "meta/meta-string"
	set "meta/metaset"
	"reflect"
	"sync"
)

type Interface interface {
	Init()
	Invoke(eventType reflect.Type, p ...Payload)
	InvokeChannel(eventType reflect.Type, channels *[]string, p ...Payload)
	Register(eventType reflect.Type, l ListenerInterface, channel ...string)
	Unregister(eventType reflect.Type, l ListenerInterface, channel ...string)
}

type Event struct {
	mutex            sync.RWMutex
	listeners        set.Set[ListenerInterface]
	channelListeners map[string]set.Set[ListenerInterface]
}

func (e *Event) Init() {
	e.listeners = set.New[ListenerInterface]()
	e.channelListeners = make(map[string]set.Set[ListenerInterface])
}

func (e *Event) Invoke(eventType reflect.Type, p ...Payload) {
	e.InvokeChannel(eventType, nil, p...)
}

func (e *Event) InvokeChannel(eventType reflect.Type, channels *[]string, p ...Payload) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	ListenerSet := set.New[ListenerInterface]()
	for l := range e.listeners {
		ListenerSet.Add(l)
	}

	if channels != nil {
		for _, channel := range *channels {
			if oldSet, ok := e.channelListeners[channel]; ok {
				ListenerSet = ListenerSet.Union(oldSet)
			}
		}
	}

	finalListeners := ListenerSet.ToSortSlice(
		func(a, b ListenerInterface) bool {
			return true
		},
	)

	if metaflag.IsDebugEvent() {
		slog.Info(
			"[Event] Invoke",
			"event", eventType.String(),
			"channels", channels,
			"listeners", metastring.ConvertStringObjects(10, finalListeners...),
			"payload", metaformat.FormatByJson(p),
		)
	}

	for _, l := range finalListeners {
		l.OnEventInvoked(eventType, p...)
	}
}

func (e *Event) Register(eventType reflect.Type, l ListenerInterface, channel ...string) {
	e.mutex.Lock()
	if len(channel) == 0 {
		e.listeners.Add(l)
	} else {
		for _, v := range channel {
			if oldSet, ok := e.channelListeners[v]; !ok {
				newSet := set.New[ListenerInterface]()
				newSet.Add(l)
				e.channelListeners[v] = newSet
			} else {
				oldSet.Add(l)
			}
		}
	}
	e.mutex.Unlock()

	if metaflag.IsDebugEvent() {
		slog.Info(
			"[Event] Register listener",
			"event", eventType.String(),
			"listener", l.GetName(),
			"channel", channel,
		)
	}
}

func (e *Event) Unregister(eventType reflect.Type, l ListenerInterface, channel ...string) {
	e.mutex.Lock()
	if len(channel) == 0 {
		e.listeners.Remove(l)
	} else {
		for _, v := range channel {
			if oldSet, ok := e.channelListeners[v]; ok {
				oldSet.Remove(l)
			}
		}
	}
	e.mutex.Unlock()

	if metaflag.IsDebugEvent() {
		slog.Info(
			"[Event] Unregister listener",
			"event", eventType.String(),
			"listener", l.GetName(),
			"channel", channel,
		)
	}
}
