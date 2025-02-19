package parser_test

import (
	"cdecl-lsp/parser"
	"slices"
	"testing"
)

func Check(t *testing.T, expected, actual []string) {
	if !slices.Equal(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestTokenize(t *testing.T) {
	s := "     here are some words    asd   "
	expected := []string{"here", "are", "some", "words", "asd"}
	tok := parser.Tokenize(s)
	Check(t, expected, tok)

	s = "word"
	expected = []string{"word"}
	tok = parser.Tokenize(s)
	Check(t, expected, tok)

	s = ""
	expected = []string{""}
	tok = parser.Tokenize(s)
	Check(t, expected, tok)
}

func TestIsDeclaration(t *testing.T) {
	s := "int foo;"
	expected := true
	actual := parser.IsDeclaration(s)
	if expected != actual {
		t.Fatalf("expected %t, got %t", expected, actual)
	}
}
