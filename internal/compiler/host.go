package compiler

import "github.com/demouth/mini-typescript-go/internal/vfs"

type CompilerHost interface {
	FS() vfs.FS
	GetCurrentDirectory() string
}

type compilerHost struct {
	currentDirectory string
	fs               vfs.FS
}

func NewCompilerHost(currentDirectory string, fs vfs.FS) CompilerHost {
	h := &compilerHost{}
	h.currentDirectory = currentDirectory
	h.fs = fs
	return h
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.currentDirectory
}

func (h *compilerHost) FS() vfs.FS {
	return h.fs
}
