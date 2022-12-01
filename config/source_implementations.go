package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

var sourceEngines []SourceEngineLoader

type SourceEngineLoader interface {
	CanParse(config *yaml.Node) bool
	GetEngine(config *yaml.Node) (SourceEngine, error)
}

func RegisterSourceEngine(engine SourceEngineLoader) {
	sourceEngines = append(sourceEngines, engine)
}

type SourceEngine interface {
	GetImageReader() (io.ReadCloser, string, error)
	IsSensitive() bool
	Caption() (string, error)
	GetMd5Hash() string
	GetTags() []string
}
