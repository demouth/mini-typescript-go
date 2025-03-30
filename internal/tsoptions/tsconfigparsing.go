package tsoptions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/demouth/mini-typescript-go/internal/core"
	"github.com/demouth/mini-typescript-go/internal/tspath"
)

func ParseJsonSourceFileConfigFileContent(basePath string) ParsedCommandLine {
	result := parseJsonConfigFileContentWorker(basePath)
	return result
}

func parseJsonConfigFileContentWorker(basePath string) ParsedCommandLine {
	basePathForFileNames := tspath.NormalizePath(basePath)

	getFileNames := func(basePath string) []string {
		fileNames := getFileNamesFromConfigSpecs(basePath)
		return fileNames
	}

	return ParsedCommandLine{
		ParsedConfig: &core.ParsedOptions{
			FileNames: getFileNames(basePathForFileNames),
		},
	}
}

func getFileNamesFromConfigSpecs(basePath string) []string {
	wildcardFileMap := make([]string, 0)
	{
		files := readDirectory(basePath)
		wildcardFileMap = append(wildcardFileMap, files...)
	}

	files := make([]string, 0, len(wildcardFileMap))
	files = append(files, wildcardFileMap...)
	return files
}

func readDirectory(currentDir string) []string {
	// TODO: Implement this function

	files := make([]string, 0)
	filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".ts") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}
