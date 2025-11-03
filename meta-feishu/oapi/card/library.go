package card

import (
	"encoding/json"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func MessageCardElementJson(e larkcard.MessageCardElement) ([]byte, error) {
	data, err := larkcore.StructToMap(e)
	if err != nil {
		return nil, err
	}
	data["tag"] = e.Tag()
	return json.Marshal(data)
}

func MessageCardBehaviorJson(e MessageCardBehavior) ([]byte, error) {
	data, err := larkcore.StructToMap(e)
	if err != nil {
		return nil, err
	}
	data["type"] = e.Type()
	return json.Marshal(data)
}
