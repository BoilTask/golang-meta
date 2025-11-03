package metaslice

func GetLastNElements[T any](slice []T, n int) []T {
	if len(slice) < n {
		return slice // 如果切片长度不足 n，直接返回整个切片
	}
	return slice[len(slice)-n:]
}

func RemoveAll[T comparable](s []T, val T) []T {
	result := s[:0] // 使用原始 slice 的底层数组，避免内存分配
	for _, v := range s {
		if v != val {
			result = append(result, v)
		}
	}
	return result
}

func RemoveAllFunc[T any](s []T, shouldRemove func(T) bool) []T {
	result := s[:0] // 原地过滤，避免额外分配
	for _, v := range s {
		if !shouldRemove(v) {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicate[T comparable](s []T) []T {
	if len(s) <= 1 {
		return s // 如果切片长度小于等于1，直接返回
	}
	seen := make(map[T]struct{})
	result := make([]T, 0, len(s)) // 预分配空间
	for _, v := range s {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
