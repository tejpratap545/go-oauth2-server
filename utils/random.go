package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math"
)

func RandomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}

func GenerateRandomSecret(n int) string {
	buff := make([]byte, n)
	rand.Read(buff)
	return base64.URLEncoding.EncodeToString(buff)
}
