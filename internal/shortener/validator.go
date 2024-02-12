package shortener

import (
	"net/url"
	"unicode"
)

func validateShortURL(value string) bool {
	for _, c := range value {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

func validateTargetURL(value string) bool {
	if _, err := url.ParseRequestURI(value); err == nil {
		return true
	}
	return false
}
