package utils

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"sync"
)

const (
	// 字符集
	AlphaCharset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumericCharset  = "0123456789"
	SpecialCharset  = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	AlphaNumCharset = AlphaCharset + NumericCharset
	AllCharset      = AlphaCharset + NumericCharset + SpecialCharset
)

var (
	// 使用加密安全的随机数生成器
	seededRand = rand.Reader
	mu         sync.Mutex
)

// RandomString 生成指定长度的随机字符串
func RandomString(length int, charset string) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	for i := range b {
		// 使用加密安全的随机数生成器
		n, err := rand.Int(seededRand, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err) // 随机数生成失败是严重错误
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// RandomInt 生成指定范围内的随机整数
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}

	// 使用加密安全的随机数生成器
	n, err := rand.Int(seededRand, big.NewInt(int64(max-min+1)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64()) + min
}

// RandomBytes 生成指定长度的随机字节
func RandomBytes(length int) []byte {
	if length <= 0 {
		return nil
	}

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// RandomUUID 生成UUID
func RandomUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// 设置UUID版本和变体
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant is 10

	return string(b)
}

// RandomFloat64 生成0到1之间的随机浮点数
func RandomFloat64() float64 {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// 将字节转换为uint64
	u := binary.BigEndian.Uint64(b)

	// 转换为0到1之间的浮点数
	return float64(u) / float64(1<<64)
}

// Shuffle 随机打乱切片
func Shuffle[T any](arr []T) {
	if len(arr) <= 1 {
		return
	}

	// Fisher-Yates shuffle算法
	for i := len(arr) - 1; i > 0; i-- {
		j := RandomInt(0, i)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
