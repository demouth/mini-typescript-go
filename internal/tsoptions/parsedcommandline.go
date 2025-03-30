package tsoptions

import "github.com/demouth/mini-typescript-go/internal/core"

type ParsedCommandLine struct {
	ParsedConfig *core.ParsedOptions
}

func (p *ParsedCommandLine) FileNames() []string {
	return p.ParsedConfig.FileNames
}
