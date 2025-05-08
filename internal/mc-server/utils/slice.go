package utils

import (
	"sort"
)

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

// Unique 返回去重后的切片
func Unique[T comparable](arr []T) []T {
	if len(arr) == 0 {
		return arr
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(arr))

	for _, v := range arr {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Filter 根据条件过滤切片
func Filter[T any](arr []T, predicate func(T) bool) []T {
	if len(arr) == 0 {
		return arr
	}

	result := make([]T, 0, len(arr))
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}

// Map 对切片中的每个元素进行转换
func Map[T, U any](arr []T, transform func(T) U) []U {
	if len(arr) == 0 {
		return []U{}
	}

	result := make([]U, len(arr))
	for i, v := range arr {
		result[i] = transform(v)
	}

	return result
}

// Sort 对切片进行排序
func Sort[T any](arr []T, less func(i, j int) bool) {
	if len(arr) <= 1 {
		return
	}

	sort.Slice(arr, less)
}

// Equal 检查两个切片是否相等
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Intersect 返回两个切片的交集
func Intersect[T comparable](a, b []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return []T{}
	}

	set := make(map[T]struct{})
	for _, v := range a {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range b {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}

	return result
}

// Difference 返回两个切片的差集
func Difference[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return []T{}
	}
	if len(b) == 0 {
		return a
	}

	set := make(map[T]struct{})
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)
	for _, v := range a {
		if _, ok := set[v]; !ok {
			result = append(result, v)
		}
	}

	return result
}
