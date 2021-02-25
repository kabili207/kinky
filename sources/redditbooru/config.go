package redditbooru

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

type SourceConfig struct {
	Filter []Filter
}

func (sc *SourceConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return ErrInvalidConfig
	}

	for _, node := range value.Content {
		val := node.Value
		var f Filter

		num, err := strconv.Atoi(val)
		if err != nil {
			var found bool
			f, found = FilterMap[val]
			if !found {
				return ErrInvalidConfig
			}
		} else {
			f = Filter(num)
		}

		sc.Filter = append(sc.Filter, f)
	}

	return nil
}
