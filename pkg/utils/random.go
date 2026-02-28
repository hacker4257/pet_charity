package utils

import (
	crypto_rand "crypto/rand"
	"math/rand"
	"time"
)

//生成验证码（安全随机）
func RandomCode(n int) string {
	code := make([]byte, n)
	buf := make([]byte, n)
	if _, err := crypto_rand.Read(buf); err == nil {
		for i := range code {
			code[i] = '0' + buf[i]%10
		}
		return string(code)
	}
	// fallback
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range code {
		code[i] = '0' + byte(r.Intn(10))
	}
	return string(code)
}
