package main

import (
	"os"
	"regexp"
	"strings"
)

type TokenType int

const (
	EndOfFile   TokenType = iota
	ClassName             = iota
	Requirement           = iota
	MethodName            = iota
	IgnoredLine           = iota
)

type TokenPattern struct {
	Type    TokenType
	Pattern *regexp.Regexp
}

type Token struct {
	Type  TokenType
	Value string
	Filename string
	Line int
}

var Patterns []TokenPattern = []TokenPattern{
	// C-Sharp matchers
	{
		Type:    ClassName,
		Pattern: regexp.MustCompile(`public\s*class\s+([A-Za-z0-9_]+)`),
	},
	{
		Type:    Requirement,
		Pattern: regexp.MustCompile(`//+\s+Requirement\s+([A-Za-z0-9-\.]+)`),
	},
	{
		Type:    MethodName,
		Pattern: regexp.MustCompile(`\s*public\s+void\s+([A-Za-z0-9_]+)\s*\(`),
	},
}

// GetToken returns the token type found on the given line.
// It is limited to a single token per line.
func GetToken(line string, patterns []TokenPattern) Token {

	for _, p := range Patterns {
		a := p.Pattern.FindStringSubmatch(line)
		if a != nil {
			return Token {
				Type: p.Type,
				Value: a[1],
			}
		}
	}

	return Token{
		Type:  IgnoredLine,
		Value: "",
	}
}

// TokenizeStream returns a list of all the tokens in a given file.
func TokenizeStream(bytes []byte, patterns []TokenPattern) []Token {

	text := string(bytes)
	lines := strings.Split(text, "\n")

	var tokens []Token

	for i, line := range(lines) {
		t := GetToken(line, patterns)
		if t.Type != IgnoredLine {
			t.Line = i+1
			tokens = append(tokens, t)
		}
	}

	return tokens
}

// TokenizeFile parses the given file for known tokens.
// It will return an array of all of the tokens found in the file
// Tokens will be populated with the filename relative to the top of
// the repository.
func TokenizeFile(filename string) ([]Token, error) {

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []Token{}, err
	}

	tokens := TokenizeStream(bytes, Patterns)

	return tokens, nil
}
