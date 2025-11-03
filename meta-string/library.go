package metastring

import (
	"encoding/json"
	"fmt"
	"io"
	metaerror "meta/meta-error"
	"os"
	"strings"
)

func IsEmpty(s *string) bool {
	return s == nil || *s == ""
}

func RemoveEmpty(s []string) []string {
	var result []string
	for _, v := range s {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicate(s []string) []string {
	if len(s) <= 1 {
		return s
	}
	seen := make(map[string]struct{})
	var result []string
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func Join[T any](s []T, sep string) string {
	var result strings.Builder
	for i, v := range s {
		if i > 0 {
			result.WriteString(sep)
		}
		result.WriteString(fmt.Sprint(v))
	}
	return result.String()
}

func JsonWithoutError(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

// GetTextEllipsis 获取文本的省略形式
func GetStringEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// GetTextEllipsis 获取文本的省略形式，会保证汉字等字符不被截断，一个字符算一个长度
func GetTextEllipsis(s string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	// 将字符串转换为 rune 切片
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}
	return string(runes[:maxLength]) + "..."
}

func getFallbackString(fallback ...string) string {
	if len(fallback) > 1 {
		return fmt.Sprint(fallback)
	} else if len(fallback) == 1 {
		return fallback[0]
	}
	return ""
}

// IsStringEqual 根据指针判断两个字符串是否相等
func IsStringEqual(s1, s2 *string) bool {
	if s1 == nil || s2 == nil {
		return s1 == s2
	}
	return *s1 == *s2
}

func GetStringSafe(s *string, fallback ...string) string {
	if s == nil {
		return getFallbackString(fallback...)
	}
	return *s
}

func GetStringFallback(s *string, fallback ...string) string {
	if s == nil || *s == "" {
		return getFallbackString(fallback...)
	}
	return *s
}

// GetStringForward 获取字符串，如果为空则使用后面的fallback字符串，否则使用forward字符串
func GetStringForward(s *string, forward string, fallback ...string) string {
	if s == nil || *s == "" {
		return getFallbackString(fallback...)
	}
	return forward
}

func GetStringFromOpenFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", metaerror.Wrap(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			metaerror.Wrap(err)
		}
	}(file)
	content, err := io.ReadAll(file)
	if err != nil {
		return "", metaerror.Wrap(err)
	}
	return string(content), nil
}

func WriteStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return metaerror.Wrap(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			metaerror.Wrap(err)
		}
	}(file)
	_, err = file.WriteString(content)
	if err != nil {
		return metaerror.Wrap(err)
	}
	return nil
}

func TrimSuffixAll(content string, suffix string) string {
	for strings.HasSuffix(content, suffix) {
		content = content[:len(content)-len(suffix)]
	}
	return content
}
