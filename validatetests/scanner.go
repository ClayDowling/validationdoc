package main

import (
	"container/list"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TokenType int

const (
	EndOfFile        TokenType = iota
	ClassName                  = iota
	RequirementLabel           = iota
	MethodName                 = iota
	IgnoredLine                = iota
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
	Id         string
	FileName   string
	ClassName  string
	MethodName string
	Line       int
}

func (rr *RequirementReference) FunctionName() string {
	if rr.ClassName == "" {
		return rr.MethodName
	}
	return fmt.Sprintf("%s.%s", rr.ClassName, rr.MethodName)
}

// References is a List of all RequirementReference objects found while parsing a folder.
var References = list.New()

var Patterns []TokenPattern = []TokenPattern{
	// C-Sharp matchers
	{
		Type:    ClassName,
		Pattern: regexp.MustCompile(`public\s*class\s+([A-Za-z0-9_]+)`),
	},
	{
		Type:    RequirementLabel,
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
			if os.IsNotExist(err) {
				continue
			}
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
	requirement := list.New()

	references := []RequirementReference{}

	for _, t := range tokens {
		switch t.Type {
		case ClassName:
			classname = t
		case RequirementLabel:
			requirement.PushBack(t)
		case MethodName:
			for e := requirement.Front(); e != nil; e = e.Next() {
				req := e.Value.(Token)
				r := RequirementReference{
					ClassName:  classname.Value,
					MethodName: t.Value,
					Id:         req.Value,
					FileName:   t.Filename,
					Line:       req.Line,
				}
				references = append(references, r)
			}
			requirement.Init()
		}
	}

	return references
}

// TokenizeFolder finds all RequirementLabel tokens in a folder and puts them in References.
// In the event of an error it will return error, otherwise nil.
func TokenizeFolder(foldername string) error {

	err := filepath.WalkDir(foldername, TokenizeWalkFunc)
	if err != nil {
		return err
	}

	return nil
}

func TokenizeWalkFunc(path string, d fs.DirEntry, err error) error {

	if d.Type().IsRegular() {

		tokens, err := TokenizeFile(path)
		if err != nil {
			return err
		}

		refs := ParseTokens(tokens)

		for _, r := range refs {
			References.PushBack(r)
		}
	}

	return nil
}

// Return all RequirementReferences which match the given ID in the global References list.
func ReferencesForRequirement(id string) []RequirementReference {

	var result = []RequirementReference{}

	for e := References.Front(); e != nil; e = e.Next() {
		r := e.Value.(RequirementReference)
		if r.Id == id {
			result = append(result, r)
		}
	}

	return result
}
