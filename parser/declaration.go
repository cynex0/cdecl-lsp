package parser

import (
	"regexp"
	"strings"
)

func IsIdentifier(tok string) bool {
	return false
}

func IsType(tok string) bool {
	return false
}

func IsDeclaration(line string) bool {
	re := regexp.MustCompile("(?:void|char|short|int|long|float|double)")
	return re.MatchString(line)
}

func Tokenize(line string) []string {
	re := regexp.MustCompile("\\s+")
	return re.Split(strings.TrimSpace(line), -1)
}
