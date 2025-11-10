package metastring

import (
	"fmt"
	"meta/object"
	"strconv"
	"strings"
)

func ConvertStringObjects[T object.Interface](limitCount int, objects ...T) string {
	objectSize := len(objects)
	result := fmt.Sprintf("[%d] {", objectSize)
	for i, obj := range objects {
		if i > 0 {
			result += ","
		}
		if limitCount >= 0 && i >= limitCount {
			result += "..."
			break
		}
		result += obj.GetName()
	}
	result += "}"
	return result
}

func Atoi64(s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return i64, err
}

func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func LowerSlice(src []string) []string {
	dst := make([]string, len(src))
	for i, s := range src {
		dst[i] = strings.ToLower(s)
	}
	return dst
}
