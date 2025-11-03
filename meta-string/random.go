package metastring

import (
	"meta/meta-math"
	"strings"
)

func GetRandomString(length int) string {
	return GetRandomStringWithParams(length, length,
		true, true, true, false, false,
	)
}

func GetRandomStringWithParams(
	minLength int,
	maxLength int,
	includeNumber, includeLower, includeUpper, includeSpecial, ignoreSimilar bool,
) string {
	var strBuilder strings.Builder

	// 根据选项构建字符集
	if includeNumber {
		if ignoreSimilar {
			strBuilder.WriteString("23456789")
		} else {
			strBuilder.WriteString("0123456789")
		}
	}
	if includeLower {
		if ignoreSimilar {
			strBuilder.WriteString("abcdefghjkmnpqrstuvwxyz")
		} else {
			strBuilder.WriteString("abcdefghijklmnopqrstuvwxyz")
		}
	}
	if includeUpper {
		if ignoreSimilar {
			strBuilder.WriteString("ABCDEFGHJKMNPQRSTUVWXYZ")
		} else {
			strBuilder.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
	}
	if includeSpecial {
		strBuilder.WriteString("!@#$%^&*()_+")
	}

	// 获取字符集
	charSet := strBuilder.String()
	if len(charSet) == 0 {
		return "" // 如果没有可用字符，返回空字符串
	}

	length := metamath.GetRandomInt(minLength, maxLength)

	// 随机生成指定长度的字符串
	var result strings.Builder
	for i := 0; i < length; i++ {
		index := metamath.GetRandomInt(0, len(charSet)-1)
		result.WriteByte(charSet[index])
	}

	return result.String()
}
