package metaformat

import (
	"encoding/json"
	"fmt"
)

func Format(format ...any) string {
	if len(format) == 0 {
		return ""
	}
	if len(format) == 1 {
		return fmt.Sprint(format[0])
	}
	formatString, ok := format[0].(string)
	if !ok {
		return fmt.Sprint(format)
	}
	return fmt.Sprintf(formatString, format[1:]...)
}

func FormatByJson(format ...any) string {
	if len(format) == 0 {
		return ""
	}
	if len(format) == 1 {
		b, err := json.Marshal(format[0])
		if err == nil {
			return string(b)
		}
		return fmt.Sprint(format[0])
	}
	formatString, ok := format[0].(string)
	if ok {
		var realFormatParam []any
		for i := 1; i < len(format); i++ {
			b, err := json.Marshal(format[i])
			if err == nil {
				realFormatParam = append(realFormatParam, string(b))
			} else {
				realFormatParam = append(realFormatParam, format[i])
			}
		}
		return fmt.Sprintf(formatString, realFormatParam...)
	} else {
		var realFormatParam []any
		for i := 0; i < len(format); i++ {
			b, err := json.Marshal(format[i])
			if err == nil {
				realFormatParam = append(realFormatParam, string(b))
			} else {
				realFormatParam = append(realFormatParam, format[i])
			}
		}
		return fmt.Sprint(realFormatParam)
	}
}

func StringByJson[T any](format ...T) string {
	if len(format) == 0 {
		return ""
	}
	if len(format) == 1 {
		b, err := json.Marshal(format[0])
		if err == nil {
			return string(b)
		}
		return fmt.Sprint(format[0])
	}
	var realFormatParam []any
	for i := 0; i < len(format); i++ {
		b, err := json.Marshal(format[i])
		if err == nil {
			realFormatParam = append(realFormatParam, string(b))
		} else {
			realFormatParam = append(realFormatParam, format[i])
		}
	}
	return fmt.Sprint(realFormatParam)
}
