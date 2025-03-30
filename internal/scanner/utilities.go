package scanner

import (
	"unicode/utf8"

	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/stringutil"
)

func GetSourceTextOfNodeFromSourceFile(sourceFile *ast.SourceFile, node *ast.Node, includeTrivia bool) string {
	return GetTextOfNodeFromSourceText(sourceFile.Text, node, includeTrivia)
}

func GetTextOfNodeFromSourceText(sourceText string, node *ast.Node, includeTrivia bool) string {
	pos := node.Pos()
	if !includeTrivia {
		pos = SkipTrivia(sourceText, pos)
	}
	text := sourceText[pos:node.End()]
	return text
}

type SkipTriviaOptions struct {
	StopAfterLineBreak bool
	StopAtComments     bool
	InJSDoc            bool
}

func SkipTrivia(text string, pos int) int {
	return SkipTriviaEx(text, pos, nil)
}

func SkipTriviaEx(text string, pos int, options *SkipTriviaOptions) int {
	if ast.PositionIsSynthesized(pos) {
		return pos
	}
	for {
		ch, size := utf8.DecodeRuneInString(text[pos:])
		switch ch {
		default:
			if ch > rune(maxAsciiCharacter) && stringutil.IsWhiteSpaceLike(ch) {
				pos += size
				continue
			}
		}
		return pos
	}
}

var (
	maxAsciiCharacter byte = 127
)
