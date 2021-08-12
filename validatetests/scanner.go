package main

import "regexp"

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
}

var Patterns []TokenPattern = []TokenPattern{
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
		Type:  EndOfFile,
		Value: "Nada",
	}
}
