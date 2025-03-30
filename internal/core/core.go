package core

import (
	"slices"
	"unicode/utf8"

	"github.com/demouth/mini-typescript-go/internal/stringutil"
)

func Filter[T any](slice []T, f func(T) bool) []T {
	for i, value := range slice {
		if !f(value) {
			result := slices.Clone(slice[:i])
			for i++; i < len(slice); i++ {
				value = slice[i]
				if f(value) {
					result = append(result, value)
				}
			}
			return result
		}
	}
	return slice
}

func FindIndex[T any](slice []T, f func(T) bool) int {
	for i, value := range slice {
		if f(value) {
			return i
		}
	}
	return -1
}

func IfElse[T any](b bool, whenTrue T, whenFalse T) T {
	if b {
		return whenTrue
	}
	return whenFalse
}

func ComputeLineStarts(text string) []TextPos {
	var result []TextPos
	pos := 0
	lineStart := 0
	for pos < len(text) {
		b := text[pos]
		if b < 0x7F {
			pos++
			switch b {
			case '\r':
				if pos < len(text) && text[pos] == '\n' {
					pos++
				}
				fallthrough
			case '\n':
				result = append(result, TextPos(lineStart))
				lineStart = pos
			}
		} else {
			ch, size := utf8.DecodeRuneInString(text[pos:])
			pos += size
			if stringutil.IsLineBreak(ch) {
				result = append(result, TextPos(lineStart))
				lineStart = pos
			}
		}
	}
	result = append(result, TextPos(lineStart))
	return result
}
