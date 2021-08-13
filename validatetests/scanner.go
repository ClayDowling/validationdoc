package main

import (
	"os"
	"path/filepath"
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
	Type     TokenType
	Value    string
	Filename string
	Line     int
}

type RequirementReference struct {
	Requirement string
	FileName    string
	ClassName   string
	MethodName  string
	Line        int
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
			return Token{
				Type:  p.Type,
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

	for i, line := range lines {
		t := GetToken(line, patterns)
		if t.Type != IgnoredLine {
			t.Line = i + 1
			tokens = append(tokens, t)
		}
	}

	return tokens
}

func gitRoot(filename string) (string, error) {

	parts := strings.Split(filename, string(os.PathSeparator))
	for i := len(parts) - 2; i >= 0; i-- {
		path := strings.Join(parts[0:i], string(os.PathSeparator))
		candidate := path + string(os.PathSeparator) + ".git"
		s, err := os.Stat(candidate)
		if err != nil {
			return "", err
		}
		if s.IsDir() {
			return path, nil
		}
	}

	return "", nil
}

// TokenizeFile parses the given file for known tokens.
// It will return an array of all of the tokens found in the file
// Tokens will be populated with the filename relative to the top of
// the repository.
func TokenizeFile(filename string) ([]Token, error) {

	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return []Token{}, err
	}

	repo, err := gitRoot(fullpath)
	if err != nil {
		return []Token{}, err
	}
	relpath, err := filepath.Rel(repo, fullpath)
	if err != nil {
		return []Token{}, err
	}
	relpath = filepath.ToSlash(relpath)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []Token{}, err
	}

	tokens := TokenizeStream(bytes, Patterns)
	for i := 0; i < len(tokens); i++ {
		tokens[i].Filename = relpath
	}

	return tokens, nil
}

// ParseTokens converts an array of tokens into an array of RequirementReferences.
func ParseTokens(tokens []Token) []RequirementReference {

	var classname Token
	var requirement Token

	references := []RequirementReference{}

	for _, t := range tokens {
		switch t.Type {
		case ClassName:
			classname = t
			break
		case Requirement:
			requirement = t
			break
		case MethodName:
			r := RequirementReference{
				ClassName:   classname.Value,
				MethodName:  t.Value,
				Requirement: requirement.Value,
				FileName:    t.Filename,
				Line:        requirement.Line,
			}
			references = append(references, r)
			break
		}
	}

	return references
}
