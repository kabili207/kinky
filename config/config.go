package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config parses the config file
type Config struct {
	Instance     string
	ClientID     string
	ClientSecret string
	AccessToken  string

	PostOptions struct {
		// "private", "unlisted" or "public"
		Visibility        string
		Content           string
		NSFW              bool
		AppendPostContent bool
	}

	Source struct {
		Extensions        []string
		Folder            string
		Recursive         bool
		EnableNSFWSuffix  bool
		EnableContentText bool
	}

	file string
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
