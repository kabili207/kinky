package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrNoEngine = errors.New("no compatible engine found")

// Config parses the config file
type Config struct {
	Instance    string
	AccessToken string
	RunInterval int64

	PostOptions struct {
		// "private", "unlisted" or "public"
		Visibility        string
		Content           string
		NSFW              bool
		AppendPostContent bool   `yaml:"append_post_content"`
		FolderID          string `yaml:"folder_id"`
	} `yaml:"post_options"`

	Source yaml.Node

	file string
}

func (c *Config) ParseSource() (SourceEngine, error) {
	for _, engine := range sourceEngines {
		if engine.CanParse(&c.Source) {
			return engine.GetEngine(&c.Source)
		}
	}

	return nil, ErrNoEngine
}

// Load the configuration
func (c *Config) Load(file string) error {
	if file == "" {
		file = ApplicationName + ".yml"
	}

	fs, err := os.Open(file)
	if err != nil {
		return err
	}

	c.file = file

	return yaml.NewDecoder(fs).Decode(c)
}

// Save the current config to disk
func (c *Config) Save() error {
	f, err := os.OpenFile(c.file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(c)
}

// Save the current config to disk
func (c *Config) SaveTo(file string) error {
	c.file = file

	return c.Save()
}
