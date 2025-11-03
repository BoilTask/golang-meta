package string

import "strings"

func GetStringByStringSlice(stringSlice []string) string {
	return strings.Join(stringSlice, "\n")
}
