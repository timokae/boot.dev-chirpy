package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			input:    "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			input:    "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			input:    "I really need a kerfuffle to go to bed sooner, Fornax!",
			expected: "I really need a **** to go to bed sooner, Fornax!",
		},
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	for _, c := range cases {
		actual := getCleanedBody(c.input, badWords)
		if actual != c.expected {
			t.Errorf("Cleaned message don't match\nActual:   '%v'\nExpected: '%v'\n\n", actual, c.expected)
		}
	}
}
