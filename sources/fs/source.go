package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/averagesecurityguy/random"
	"github.com/bmatcuk/doublestar/v3"
)

type Source struct {
	config *SourceConfig
	file   *string
}

func New(cfg *SourceConfig) (*Source, error) {
	stat, err := os.Stat(cfg.Folder)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, ErrNotFolder
	}

	if cfg.Extensions == nil || len(cfg.Extensions) == 0 {
		return nil, fmt.Errorf("file extension list empty: %w", ErrInvalidConfig)
	}

	s := new(Source)
	s.config = cfg

	return s, nil
}

func (s *Source) getFiles() ([]string, error) {
	var files []string
	var err error

	if s.config.Recursive {
		files, err = doublestar.Glob(path.Join(s.config.Folder, "**"))
	} else {
		files, err = doublestar.Glob(path.Join(s.config.Folder, "*"))
	}

	if err != nil {
		return nil, err
	}

	var filteredFiles []string
	startIndex := 0
	for idx, el := range files {
		s, err := os.Stat(el)
		if err != nil {
			continue
		}

		if s.IsDir() {
			filteredFiles = append(filteredFiles, files[startIndex:idx]...)
			startIndex++
		}
	}
	filteredFiles = append(filteredFiles, files[startIndex:]...)

	return filteredFiles, nil
}

func (s *Source) getFile() (string, error) {
	if s.file == nil {
		files, err := s.getFiles()
		if err != nil {
			return "", err
		}

		index, err := random.Uint64Range(0, uint64(len(files)-1))
		if err != nil {
			return "", err
		}

		fil, err := filepath.Abs(files[index])
		if err != nil {
			return "", err
		}

		s.file = &fil
	}

	return *s.file, nil
}

func (s *Source) Caption() (string, error) {
	if !s.config.EnableContentText {
		return "", nil
	}

	fil, err := s.getFile()
	if err != nil {
		return "", err
	}

	fileDir := filepath.Dir(fil)
	fileName := filepath.Base(fil)
	fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".txt"
	contentFile := filepath.Join(fileDir, fileName)

	_, err = os.Stat(contentFile)
	if err != nil {
		return "", nil
	}

	btxt, err := ioutil.ReadFile(contentFile)
	if err != nil {
		return "", err
	}

	return string(btxt), nil
}

func (s *Source) GetImageReader() (io.ReadCloser, error) {
	fil, err := s.getFile()
	if err != nil {
		return nil, err
	}

	return os.Open(fil)
}

func (s *Source) IsSensitive() bool {
	fil, err := s.getFile()
	if err != nil {
		return false
	}

	if s.config.EnableNSFWSuffix {
		fnameParts := strings.Split(filepath.Base(fil), ".")
		if len(fnameParts) >= 3 && fnameParts[len(fnameParts)-2] == "nsfw" {
			return true
		}
	}

	if s.config.EnableNSFWFolder {
		parentFolder := filepath.Base(filepath.Dir(fil))
		if parentFolder == "nsfw" {
			return true
		}
	}

	return false
}
