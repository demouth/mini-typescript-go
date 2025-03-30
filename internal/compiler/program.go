package compiler

import (
	"sync"

	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/printer"
	"github.com/demouth/mini-typescript-go/internal/tsoptions"
)

type ProgramOptions struct {
	Host CompilerHost
}
type Program struct {
	host  CompilerHost
	files []*ast.SourceFile
}

type EmitResult struct {
}

func NewProgram(options ProgramOptions) *Program {
	p := &Program{}
	p.host = options.Host

	parseConfigFileContent := tsoptions.ParseJsonSourceFileConfigFileContent(
		p.host.GetCurrentDirectory(),
	)

	rootFiles := parseConfigFileContent.FileNames()

	p.files = processAllProgramFiles(p.host, rootFiles)
	return p
}
func (p *Program) SourceFiles() []*ast.SourceFile {
	return p.files
}
func (p *Program) Emit() *EmitResult {

	host := &emitHost{
		program: p,
	}

	writerPool := &sync.Pool{
		New: func() any {
			return printer.NewTextWriter("\n")
		},
	}

	sourceFiles := getSourceFilesToEmit(host)
	for _, sourceFile := range sourceFiles {
		emitter := &emitter{
			host:       host,
			sourceFile: sourceFile,
		}
		writer := writerPool.Get().(printer.EmitTextWriter)
		writer.Clear()
		emitter.writer = writer
		emitter.paths = getOutputPathsFor(sourceFile, host)
		emitter.emit()
		emitter.writer = nil
		writerPool.Put(writer)
	}
	return &EmitResult{}
}
