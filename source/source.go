package source

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattn/go-mastodon"

	"z0ne.dev/kura/kinky/config"
)

type UploadFunc func(string) (mastodon.ID, error)

type Source struct {
	config *config.Config
}

func New(cfg *config.Config) (*Source, error) {
	stat, err := os.Stat(cfg.Source.Folder)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, ErrNotFolder
	}

	if cfg.Source.Extensions == nil || len(cfg.Source.Extensions) == 0 {
		return nil, fmt.Errorf("file extension list empty: %w", ErrInvalidConfig)
	}

	s := new(Source)
	s.config = cfg

	return s, nil
}

func listContains(list []string, entry string) bool {
	for _, e := range list {
		if e == entry {
			return true
		}
	}

	return false
}

func (s *Source) getFiles() []string {
	var files []string

	var glob func(string)
	glob = func(folder string) {
		//nolint:errcheck
		globs, _ := filepath.Glob(path.Join(folder, "*"))
		for _, f := range globs {
			//nolint:errcheck
			stat, _ := os.Stat(f)
			if stat.IsDir() && s.config.Source.Recursive {
				glob(f)
			} else if listContains(s.config.Source.Extensions, filepath.Ext(f)[1:]) {
				files = append(files, f)
			}
		}
	}
	glob(s.config.Source.Folder)

	return files
}

func (s *Source) pickFile() (string, error) {
	files := s.getFiles()

	//nolint:gosec
	calls := rand.Intn(10)
	index := 0
	for i := 0; i < calls; i++ {
		//nolint:gosec
		index += rand.Intn(len(files))
	}

	return files[index%len(files)], nil
}

func (s *Source) nsfwMagix(file string, toot *mastodon.Toot) {
	fnameParts := strings.Split(filepath.Base(file), ".")
	if s.config.Source.EnableNSFWSuffix {
		if len(fnameParts) >= 3 && fnameParts[len(fnameParts)-2] == "nsfw" {
			toot.Sensitive = true
		}
	}
}

func (s *Source) imageText(file string, toot *mastodon.Toot) error {
	if !s.config.Source.EnableContentText {
		return nil
	}

	fileDir := filepath.Dir(file)
	fileName := filepath.Base(file)
	fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".txt"
	contentFile := filepath.Join(fileDir, fileName)

	_, err := os.Stat(contentFile)
	if err != nil {
		return nil
	}
	btxt, err := ioutil.ReadFile(contentFile)
	if err != nil {
		return err
	}

	txt := string(btxt)
	if len(txt) > 0 {
		if s.config.PostOptions.AppendPostContent {
			toot.Status += txt
		} else {
			toot.Status = txt
		}
	}

	return nil
}

func (s *Source) Post(upload UploadFunc) (*mastodon.Toot, error) {
	file, err := s.pickFile()
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	toot := &mastodon.Toot{
		Visibility: s.config.PostOptions.Visibility,
		Sensitive:  s.config.PostOptions.NSFW,
		Status:     s.config.PostOptions.Content,
	}

	mediaID, err := upload(abs)
	if err != nil {
		return nil, err
	}
	toot.MediaIDs = []mastodon.ID{mediaID}

	s.nsfwMagix(file, toot)
	if err := s.imageText(file, toot); err != nil {
		return nil, err
	}

	return toot, nil
}