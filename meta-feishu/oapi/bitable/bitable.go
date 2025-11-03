package bitable

import (
	"encoding/json"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
	metastring "meta/meta-string"
	"strconv"
	"strings"
)

func getBitableStringField(v interface{}) *string {
	if v == nil {
		return nil
	}
	if v, ok := v.(string); ok {
		return &v
	}
	jsonsStr, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var filedText FiledText
	err = json.Unmarshal(jsonsStr, &filedText)
	if err != nil {
		return nil
	}
	return &filedText.Text
}

func GetBitableStringField(a *larkbitable.AppTableRecord, key string) *string {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	if v, ok := v.([]interface{}); ok {
		var finalRes string
		for _, item := range v {
			res := getBitableStringField(item)
			if res != nil {
				finalRes = finalRes + *res
			}
		}
		return &finalRes
	}
	return getBitableStringField(v)
}

func GetBitableStringsField(a *larkbitable.AppTableRecord, key string) []string {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	if v, ok := v.([]interface{}); ok {
		var finalRes []string
		for _, item := range v {
			res := getBitableStringField(item)
			if res != nil {
				finalRes = append(finalRes, *res)
			}
		}
		return finalRes
	}
	res := getBitableStringField(v)
	if res != nil {
		return []string{*res}
	}
	return nil
}

func GetBitableBoolField(a *larkbitable.AppTableRecord, key string) *bool {
	if a == nil || a.Fields == nil {
		return nil
	}

	v := a.Fields[key]
	if v == nil {
		return nil
	}

	if v, ok := v.(bool); ok {
		return &v
	}

	if v, ok := v.(string); ok {
		vLower := strings.ToLower(v)
		b := vLower == "true"
		return &b
	}

	vInt := GetBitableIntField(a, key)
	if vInt != nil {
		b := *vInt != 0
		return &b
	}

	return nil
}

func GetBitableIntField(a *larkbitable.AppTableRecord, key string) *int {
	if a == nil || a.Fields == nil {
		return nil
	}

	v := a.Fields[key]
	if v == nil {
		return nil
	}

	if v, ok := v.(int); ok {
		return &v
	}

	if v, ok := v.(string); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return &i
		}
	}

	return nil
}

func GetBitableInt64Field(a *larkbitable.AppTableRecord, key string) *int64 {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	if v, ok := v.(int64); ok {
		return &v
	}
	if v, ok := v.(int); ok {
		vv := int64(v)
		return &vv
	}
	if v, ok := v.(float64); ok {
		vv := int64(v)
		return &vv
	}
	if v, ok := v.(string); ok {
		if i, err := metastring.Atoi64(v); err == nil {
			return &i
		}
	}
	return nil
}

func GetBitableFloat64Field(a *larkbitable.AppTableRecord, key string) *float64 {
	if a == nil || a.Fields == nil {
		return nil
	}

	v := a.Fields[key]
	if v == nil {
		return nil
	}

	if v, ok := v.(float64); ok {
		return &v
	}

	return nil
}

func GetBitableUrlField(a *larkbitable.AppTableRecord, key string) *larkbitable.Url {
	if a == nil || a.Fields == nil {
		return nil
	}

	v := a.Fields[key]
	if v == nil {
		return nil
	}

	if v, ok := v.(map[string]interface{}); ok {
		textInterface, okText := v["text"]
		linkInterface, okLink := v["link"]
		if okText && okLink {
			text := textInterface.(string)
			link := linkInterface.(string)
			url := larkbitable.Url{
				Text: &text,
				Link: &link,
			}
			return &url
		}
	}

	return nil
}

func GetBitablePersonsField(a *larkbitable.AppTableRecord, key string) []larkbitable.Person {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	jsonsStr, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var persons []larkbitable.Person
	err = json.Unmarshal(jsonsStr, &persons)
	if err != nil {
		return nil
	}
	return persons
}

func GetBitableGroupField(a *larkbitable.AppTableRecord, key string) []larkbitable.Group {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	jsonsStr, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var groups []larkbitable.Group
	err = json.Unmarshal(jsonsStr, &groups)
	if err != nil {
		return nil
	}
	return groups
}

func GetBitableLinkRecordField(a *larkbitable.AppTableRecord, key string) []string {
	if a == nil || a.Fields == nil {
		return nil
	}
	v := a.Fields[key]
	if v == nil {
		return nil
	}
	if v, ok := v.(map[string]interface{}); ok {
		values := v["link_record_ids"]
		if values != nil {
			if v, ok := values.([]interface{}); ok {
				var finalRes []string
				for _, item := range v {
					if item, ok := item.(string); ok {
						finalRes = append(finalRes, item)
					}
				}
				return finalRes
			}
		}
	}
	return nil
}
