package printer

import (
	"strings"

	"github.com/demouth/mini-typescript-go/internal/core"
)

type textWriter struct {
	newLine   string
	builder   strings.Builder
	lineStart bool
}

func (w *textWriter) Clear() {
	*w = textWriter{
		newLine:   w.newLine,
		lineStart: true,
	}
}

func NewTextWriter(newLine string) EmitTextWriter {
	var w textWriter
	w.newLine = newLine
	w.Clear()
	return &w
}

func (w *textWriter) String() string {
	return w.builder.String()
}

func (w *textWriter) writeText(s string) {
	if s != "" {
		w.builder.WriteString(s)
		w.updateLineCountAndPosFor(s)
	}
}

func (w *textWriter) writeLineRaw() {
	w.builder.WriteString(w.newLine)
	w.lineStart = true
}

func (w *textWriter) Write(s string) {
	w.writeText(s)
}

func (w *textWriter) WriteLine() {
	if !w.lineStart {
		w.writeLineRaw()
	}
}

func (w *textWriter) WriteStringLiteral(text string) {
	w.Write(text)
}
func (w *textWriter) WriteTrailingSemicolon(text string) {
	w.Write(text)
}
func (w *textWriter) updateLineCountAndPosFor(s string) {
	_ = core.ComputeLineStarts(s)
	w.lineStart = false
}
