package textutil

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(letters[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letters))])
	}
	return string(b)
}
