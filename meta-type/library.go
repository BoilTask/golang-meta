package metatype

func IsBoolTrue(value interface{}) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != ""
	case float64:
		return v != 0
	case []interface{}:
		return len(v) > 0
	default:
		return false
	}
}
