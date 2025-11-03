package metajson

import (
	"bytes"
	"encoding/json"
)

func GetJsonStr(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
