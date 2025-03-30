package compiler

import (
	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/core"
	"github.com/demouth/mini-typescript-go/internal/printer"
	"github.com/demouth/mini-typescript-go/internal/tspath"
)

type emitter struct {
	host       EmitHost
	writer     printer.EmitTextWriter
	paths      *outputPaths
	sourceFile *ast.SourceFile
}
type outputPaths struct {
	jsFilePath string
}

func (e *emitter) emit() {
	e.emitJsFile(e.sourceFile, e.paths.jsFilePath)
}

func (e *emitter) emitJsFile(sourceFile *ast.SourceFile, jsFilePath string) {
	printer := printer.NewPrinter()
	e.printSourceFile(jsFilePath, sourceFile, printer)
}

func (e *emitter) printSourceFile(jsFilePath string, sourceFile *ast.SourceFile, printer *printer.Printer) bool {
	printer.Write(sourceFile.AsNode(), sourceFile, e.writer)

	e.writer.WriteLine()

	text := e.writer.String()

	e.host.WriteFile(jsFilePath, text, false)

	e.writer.Clear()

	return true
}

func ownOutputFilePath(fileName string, extension string) string {
	emitOutputFilePathWithoutExtension := tspath.RemoveFileExtension(fileName)
	return emitOutputFilePathWithoutExtension + extension
}

func getOutputPathsFor(sourceFile *ast.SourceFile, host EmitHost) *outputPaths {
	ownOutputFilePath := ownOutputFilePath(sourceFile.FileName(), ".js")
	return &outputPaths{
		jsFilePath: ownOutputFilePath,
	}
}

func getSourceFilesToEmit(host EmitHost) []*ast.SourceFile {

	sourceFiles := host.SourceFiles()

	return core.Filter(sourceFiles, func(sourceFile *ast.SourceFile) bool {
		return sourceFileMayBeEmitted(sourceFile, host)
	})
}

func sourceFileMayBeEmitted(_ *ast.SourceFile, _ EmitHost) bool {
	return true
}
