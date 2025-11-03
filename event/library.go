package event

import (
	"log/slog"
	metaerror "meta/meta-error"
	metaflag "meta/meta-flag"
	metaformat "meta/meta-format"
	"reflect"
	"sync"
)

var (
	eventMap = make(map[reflect.Type]interface{})
	mutex    = &sync.RWMutex{}
)

func GetEvent(subsystemType reflect.Type) Interface {
	mutex.RLock()
	defer mutex.RUnlock()
	return GetEventUnsafe(subsystemType)
}

func GetEventUnsafe(subsystemType reflect.Type) Interface {
	event, ok := eventMap[subsystemType]
	if !ok {
		return nil
	}
	return event.(Interface)
}

func IsEvent[T any, _ interface {
	*T
	Interface
}](event reflect.Type) bool {
	return event == reflect.TypeFor[T]()
}

func RegisterListener[T any, _ interface {
	*T
	Interface
}](l ListenerInterface, channel ...string) {
	eventType := reflect.TypeFor[T]()
	event := GetEvent(eventType)
	if event == nil {
		mutex.Lock()
		var newEvent T
		eventMap[eventType] = &newEvent
		event = GetEventUnsafe(eventType)
		event.(Interface).Init()
		mutex.Unlock()
		event.(Interface).Register(eventType, l, channel...)
	} else {
		event.(Interface).Register(eventType, l, channel...)
	}
}

func UnregisterListener[T any, _ interface {
	*T
	Interface
}](l ListenerInterface) {
	eventType := reflect.TypeFor[T]()
	event := GetEvent(eventType)
	if event != nil {
		event.(Interface).Unregister(eventType, l)
	}
}

func Invoke[T any, _ interface {
	*T
	Interface
}](p ...Payload) {
	eventType := reflect.TypeFor[T]()
	event := GetEvent(eventType)
	if event == nil {
		if metaflag.IsDebugEvent() {
			slog.Info(
				"[Event] Invoke event without listener",
				"event", eventType.String(),
				"payload", metaformat.FormatByJson(p),
			)
		}
		return
	}
	event.(Interface).Invoke(eventType, p)
}

func InvokeChannel[T any, _ interface {
	*T
	Interface
}](channels *[]string, p ...Payload) {
	eventType := reflect.TypeFor[T]()
	event := GetEvent(eventType)
	if event == nil {
		if metaflag.IsDebugEvent() {
			slog.Info(
				"[Event] Invoke event with channel without listener",
				"event", eventType.String(),
				"channels", channels,
				"payload", metaformat.StringByJson(p...),
			)
		}
		return
	}
	event.(Interface).InvokeChannel(eventType, channels, p...)
}

func ParsePayloadIndex[T any, _ interface {
	*T
}](payload []Payload, index int) (*T, error) {
	if index < 0 || index >= len(payload) {
		return nil, metaerror.New("index %d out of range [%d,%d)", index, 0, len(payload))
	}
	res, ok := payload[index].(*T)
	if !ok {
		return nil, metaerror.New(
			"invalid payload type, payload:%s expect:%s",
			reflect.TypeOf(payload[0]),
			reflect.TypeFor[T](),
		)
	}
	return res, nil
}

func ParsePayload[T any, _ interface {
	*T
}](payload []Payload) (*T, error) {
	return ParsePayloadIndex[T](payload, 0)
}
