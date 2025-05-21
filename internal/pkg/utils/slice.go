package utils

// Contains 检查切片是否包含指定元素
func Contains[T comparable](arr []T, value T) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
