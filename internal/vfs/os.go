package vfs

import (
	"os"

	"github.com/demouth/mini-typescript-go/internal/tspath"
)

func FromOS() FS {
	return &osFS{}
}

type osFS struct {
}

func (vfs *osFS) writeFile(path string, content string, writeByteOrderMark bool) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if writeByteOrderMark {
		if _, err := file.WriteString("\uFEFF"); err != nil {
			return err
		}
	}

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func (vfs *osFS) ensureDirectoryExists(directoryPath string) error {
	return os.MkdirAll(directoryPath, 0o777)
}
func (vfs *osFS) WriteFile(path string, content string, writeByteOrderMark bool) error {
	_ = rootLength(path)
	if err := vfs.writeFile(path, content, writeByteOrderMark); err == nil {
		return nil
	}
	if err := vfs.ensureDirectoryExists(tspath.GetDirectoryPath(tspath.NormalizePath(path))); err != nil {
		return err
	}
	return vfs.writeFile(path, content, writeByteOrderMark)
}
