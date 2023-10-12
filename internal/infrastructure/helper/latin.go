package helper

import "unicode"

func IsLatinOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) || !unicode.Is(unicode.Latin, r) {
			return false
		}
	}
	return true
}

func IsLatinOrDigitOnly(s string) bool {
	for _, r := range s {
		if (!unicode.IsLetter(r) || !unicode.Is(unicode.Latin, r)) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
