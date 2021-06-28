package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrStartingAtDigit = errors.New("you pass string starting at digit")

var ErrContainsNumber = errors.New("your string contains a number")

func Unpack(rawStr string) (string, error) {
	var strBuilder strings.Builder
	rawStr = strings.TrimSpace(rawStr)

	if len(rawStr) == 0 {
		return "", nil
	}

	for idx, r := range rawStr {
		if unicode.IsDigit(r) {
			prevIdx := idx - 1
			if prevIdx < 0 {
				return "", ErrStartingAtDigit
			}
			prevRune := rune(rawStr[prevIdx])
			if unicode.IsDigit(prevRune) {
				return "", ErrContainsNumber
			}
			continue
		}

		nextIdx := idx + 1
		if nextIdx < len(rawStr) {
			nextRune := rune(rawStr[nextIdx])
			if unicode.IsDigit(nextRune) {
				countLetters := int(nextRune - '0')
				currentLetter := string(rawStr[idx])
				repeatedLetters := strings.Repeat(currentLetter, countLetters)
				strBuilder.WriteString(repeatedLetters)
				continue
			}
			strBuilder.WriteString(string(rawStr[idx]))
			continue
		}

		strBuilder.WriteString(string(rawStr[idx]))
	}

	return strBuilder.String(), nil
}
