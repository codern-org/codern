package generator

import (
	"math/rand"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var number = []byte("0123456789")
var alphaNumeric = append(charset, number...)

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(b)
}
