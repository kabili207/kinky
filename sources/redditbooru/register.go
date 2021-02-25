package redditbooru

import (
	"gopkg.in/yaml.v3"

	"z0ne.dev/kura/kinky/config"
)

func Register() {
	config.RegisterSourceEngine(newEngine())
}

type EngineLoader struct {
}

func (e *EngineLoader) CanParse(config *yaml.Node) bool {
	cfg := new(SourceConfig)
	err := config.Decode(cfg)

	return err == nil && len(cfg.Filter) > 0
}

func (e *EngineLoader) GetEngine(config *yaml.Node) (config.SourceEngine, error) {
	cfg := new(SourceConfig)
	if err := config.Decode(cfg); err != nil {
		return nil, err
	}

	return New(cfg)
}

func newEngine() *EngineLoader {
	return new(EngineLoader)
}
