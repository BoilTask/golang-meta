package metasql

import (
	"encoding/json"
	"strings"
	"time"
)

type MysqlTime struct {
	time.Time
}

func (t *MysqlTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

func (t *MysqlTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" {
		t.Time = time.Time{}
		return nil
	}
	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		t.Time = parsed
		return nil
	}
	if parsed, err := time.ParseInLocation("2006-01-02 15:04:05.000000", s, time.Local); err == nil {
		t.Time = parsed
		return nil
	}
	if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		t.Time = parsed
		return nil
	}
	return &time.ParseError{Layout: "multiple", Value: s}
}
