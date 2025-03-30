package compiler

import "github.com/demouth/mini-typescript-go/internal/ast"

type EmitHost interface {
	SourceFiles() []*ast.SourceFile
	WriteFile(fileName string, text string, writeByteOrderMark bool) error
}

type emitHost struct {
	program *Program
}

func (host *emitHost) SourceFiles() []*ast.SourceFile {
	return host.program.SourceFiles()
}

func (host *emitHost) WriteFile(fileName string, text string, writeByteOrderMark bool) error {
	return host.program.host.FS().WriteFile(fileName, text, writeByteOrderMark)
}
