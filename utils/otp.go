package utils

import (
	"crypto/rand"
	"os"
	"strconv"
)

const otpChars = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateOTPCode() (string, error) {
	otpLimit, err := strconv.Atoi(os.Getenv("OTP_CODE_LENGTH"))

	if err != nil {
		return "", err
	}

	buffer := make([]byte, otpLimit)
	_, err = rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < otpLimit; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
