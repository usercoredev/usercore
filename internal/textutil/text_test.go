package textutil

import (
	"testing"
	"unicode/utf8"
)

func TestRandomString(t *testing.T) {
	length := 10
	str := RandomString(length)
	if utf8.RuneCountInString(str) != length {
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
