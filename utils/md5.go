package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

// @function: MD5V
// @description: md5加密
// @param: str []byte
// @return: string
func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// GetRandSlatString - 获取随机盐
// @param {int} n
// @returns string
func GetRandSlatString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
