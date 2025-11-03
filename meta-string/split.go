package metastring

import "strings"

func RSplit(s, sep string, n int) []string {
	result := make([]string, 0)
	for i := 0; i < n-1; i++ {
		idx := strings.LastIndex(s, sep)
		if idx == -1 {
			break
		}
		result = append([]string{s[idx+len(sep):]}, result...)
		s = s[:idx]
	}
	return append([]string{s}, result...)
}
