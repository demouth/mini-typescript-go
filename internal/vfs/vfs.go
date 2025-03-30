package vfs

import (
	"fmt"

	"github.com/demouth/mini-typescript-go/internal/tspath"
)

type FS interface {
	WriteFile(path string, data string, writeByteOrderMark bool) error
}

func rootLength(p string) int {
	l := tspath.GetEncodedRootLength(p)
	if l <= 0 {
		panic(fmt.Sprintf("vfs: path %q is not absolute", p))
	}
	return l
}
