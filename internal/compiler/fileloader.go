package compiler

import (
	"os"

	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/parser"
	"github.com/demouth/mini-typescript-go/internal/tspath"
)

type fileLoader struct {
	host      CompilerHost
	rootTasks []*parseTask
}

func (p *fileLoader) parseSourceFile(fileName string) *ast.SourceFile {

	bytes, _ := os.ReadFile(fileName)
	text := string(bytes)

	sourceFile := parser.ParseSourceFile(fileName, text)
	return sourceFile
}

func (p *fileLoader) addRootTasks(files []string) {
	for _, fileName := range files {
		absPath := tspath.GetNormalizedAbsolutePath(fileName, p.host.GetCurrentDirectory())
		p.rootTasks = append(
			p.rootTasks,
			&parseTask{
				normalizedFilePath: absPath,
			},
		)
	}
}

func (p *fileLoader) startTasks(tasks []*parseTask) {
	if len(tasks) == 0 {
		return
	}

	for _, task := range tasks {
		task.start(p)
	}
}

func processAllProgramFiles(
	host CompilerHost,
	rootFiles []string,
) []*ast.SourceFile {

	loader := &fileLoader{
		host:      host,
		rootTasks: make([]*parseTask, 0, len(rootFiles)),
	}
	loader.addRootTasks(rootFiles)
	loader.startTasks(loader.rootTasks)

	files, libFiles := []*ast.SourceFile{}, []*ast.SourceFile{}

	for _, task := range loader.rootTasks {
		// file := &ast.SourceFile{}
		// file.SetFileName(task.normalizedFilePath) // fixme: this is a hack
		// file.SetPath(tspath.Path(task.normalizedFilePath))
		files = append(files, task.file)
	}

	return append(libFiles, files...)
}

type parseTask struct {
	normalizedFilePath string
	file               *ast.SourceFile
}

func (t *parseTask) start(loader *fileLoader) {
	file := loader.parseSourceFile(t.normalizedFilePath)
	t.file = file
}
