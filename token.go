package main

import "unicode"

type tokens []rune

func tokenize(s string) tokens {
	return append(tokens(s), 0)
}

func (ts *tokens) Current() token {
	return token((*ts)[0])
}

func (ts *tokens) Advance() {
	*ts = (*ts)[1:]
}

type token rune

func (t token) IsDigit() bool {
	return unicode.IsDigit(rune(t))
}

func (t token) IsIdent() bool {
	return unicode.IsLetter(rune(t))
}

func (t token) IsOneOf(ts ...token) bool {
	for _, c := range ts {
		if c == t {
			return true
		}
	}
	return false
}

func (t token) String() string {
	switch t {
	case 0:
		return "end of data"
	default:
		return string(t)
	}
}
