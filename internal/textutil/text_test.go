package textutil

import (
	"testing"
	"unicode/utf8"
)

func TestRandomString(t *testing.T) {
	length := int8(10)
	str := RandomString(length)
	if int8(utf8.RuneCountInString(str)) != length {
		t.Errorf("Expected string length of %d, but got %d", length, utf8.RuneCountInString(str))
	}

	for _, runeValue := range str {
		if !((runeValue >= 'a' && runeValue <= 'z') || (runeValue >= 'A' && runeValue <= 'Z') || (runeValue >= '0' && runeValue <= '9')) {
			t.Errorf("String contains invalid character: %q", runeValue)
		}
	}
}

func TestRandomStringUnique(t *testing.T) {
	str1 := RandomString(10)
	str2 := RandomString(10)
	if str1 == str2 {
		t.Errorf("Expected unique strings, but got identical %s and %s", str1, str2)
	}
}

func TestRandomStringLengthZero(t *testing.T) {
	str := RandomString(0)
	if len(str) != 0 {
		t.Errorf("Expected empty string for length 0, but got: %s", str)
	}
}

func TestRandomStringNegativeLength(t *testing.T) {
	str := RandomString(-1)
	if len(str) != 0 {
		t.Errorf("Expected empty string for negative length, but got: %s", str)
	}
}

func TestRandomInt(t *testing.T) {
	minInt := 1
	maxInt := 10
	randomIntTest := randomInt(minInt, maxInt)
	if randomIntTest < minInt || randomIntTest > maxInt {
		t.Errorf("Expected random int between %d and %d, but got %d", minInt, maxInt, randomIntTest)
	}
}

func TestGenerateOTPCode(t *testing.T) {
	otpCode, err := GenerateOTPCode()
	if err != nil {
		t.Errorf("Expected nil error, but got: %s", err)
	}

	if len(otpCode) != defaultOTPLength {
		t.Errorf("Expected %d length OTP code, but got %d", defaultOTPLength, len(otpCode))
	}

	for _, runeValue := range otpCode {
		if !((runeValue >= 'a' && runeValue <= 'z') || (runeValue >= 'A' && runeValue <= 'Z') || (runeValue >= '0' && runeValue <= '9')) {
			t.Errorf("OTP code contains invalid character: %q", runeValue)
		}
	}
}
