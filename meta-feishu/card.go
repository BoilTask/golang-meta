package metafeishu

import "github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"

type CardCallbackParam struct {
	Event    *callback.CardActionTriggerEvent
	Response *callback.CardActionTriggerResponse
}
