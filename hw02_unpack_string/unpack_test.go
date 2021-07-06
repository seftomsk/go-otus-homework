package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},

		// The new cases
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "d\n\b2abc", expected: "d\n\b\babc"},
		{input: " s2v", expected: "ssv"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

// The new cases

func TestUnpackStartingAtDigit(t *testing.T) {
	invalidStrings := []string{"3abc", "45"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrStartingAtDigit), "actual error %q", err)
		})
	}
}

func TestUnpackContainsNumber(t *testing.T) {
	const rawString = "aaa10b"
	_, err := Unpack(rawString)
	require.Truef(t, errors.Is(err, ErrContainsNumber), "actual error %q", err)
}
