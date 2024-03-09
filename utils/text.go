package utils

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var localRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateOTPCode() (string, error) {
	otpLength := getOTPLength()
	otpCode := RandomString(otpLength)
	return otpCode, nil
}

func RandomString(length int8) string {
	if length <= 0 {
		return ""
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[randomInt(0, len(letters))]
	}
	return string(b)
}
func randomInt(min, max int) int {
	return min + localRand.Intn(max-min)
}

const defaultOTPLength = 6

func getOTPLength() int8 {
	otpLimitStr := os.Getenv("OTP_CODE_LENGTH")
	if otpLimitStr == "" {
		log.Printf("OTP_CODE_LENGTH not set, using default value: %d", defaultOTPLength)
		return defaultOTPLength
	}

	otpLimit, err := strconv.Atoi(otpLimitStr)
	if err != nil {
		log.Printf("Invalid OTP_CODE_LENGTH value, using default: %d", defaultOTPLength)
		return defaultOTPLength
	}

	return int8(otpLimit)
}
