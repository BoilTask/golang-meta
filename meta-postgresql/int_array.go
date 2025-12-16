package metapostgresql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	metaerror "meta/meta-error"
)

type IntArray []int

// Value 实现 driver.Valuer 接口，将 Go 类型转换为数据库值
func (a IntArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	// 构建 PostgreSQL 数组格式 {1,2,3}
	result := "{"
	for i, v := range a {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", v)
	}
	result += "}"
	return result, nil
}

// Scan 实现 sql.Scanner 接口，将数据库值转换为 Go 类型
func (a *IntArray) Scan(value interface{}) error {
	// 处理nil值
	if value == nil {
		*a = IntArray{}
		return nil
	}

	// 转换为[]byte或string
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return metaerror.New("failed to scan IntArray: cannot convert %T to string", value)
	}

	// 去除首尾空格
	for len(str) > 0 && (str[0] == ' ' || str[0] == '\t' || str[0] == '\n' || str[0] == '\r') {
		str = str[1:]
	}
	for len(str) > 0 && (str[len(str)-1] == ' ' || str[len(str)-1] == '\t' || str[len(str)-1] == '\n' || str[len(str)-1] == '\r') {
		str = str[:len(str)-1]
	}

	// 处理空字符串
	if str == "" {
		*a = IntArray{}
		return nil
	}

	// 尝试解析为标准JSON数组 [1,2,3]
	if len(str) > 0 && str[0] == '[' && str[len(str)-1] == ']' {
		var result []int
		err := json.Unmarshal([]byte(str), &result)
		if err == nil {
			*a = IntArray(result)
			return nil
		}
	}

	// 尝试解析为PostgreSQL数组格式 {1,1}
	if len(str) > 0 && str[0] == '{' && str[len(str)-1] == '}' {
		// 去除花括号
		content := str[1 : len(str)-1]
		if content == "" {
			*a = IntArray{}
			return nil
		}

		// 分割数字
		var result IntArray
		numStr := ""
		for _, ch := range content {
			if ch >= '0' && ch <= '9' || ch == '-' {
				numStr += string(ch)
			} else if ch == ',' || ch == ' ' {
				if numStr != "" {
					var n int
					_, err := fmt.Sscanf(numStr, "%d", &n)
					if err == nil {
						result = append(result, n)
					}
					numStr = ""
				}
			}
		}
		// 处理最后一个数字
		if numStr != "" {
			var n int
			_, err := fmt.Sscanf(numStr, "%d", &n)
			if err == nil {
				result = append(result, n)
			}
		}
		if len(result) > 0 {
			*a = result
			return nil
		}
	}
	// 所有尝试都失败，返回详细的错误信息
	return metaerror.New("failed to parse IntArray from value: '%s'", str)
}
