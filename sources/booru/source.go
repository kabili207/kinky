package booru

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/juliangruber/go-intersect/v2"
	"z0ne.dev/kura/kinky/config"
)

type Source struct {
	config   *SourceConfig
	metadata *booruMetadata
}

func New(cfg *SourceConfig) (*Source, error) {
	s := new(Source)
	s.config = cfg

	return s, nil
}

func (s *Source) getMeta() (*booruMetadata, error) {
	if s.metadata == nil {
		post, ok := dialects[s.config.Dialect]
		if !ok {
			return nil, ErrDialectNotSupported
		}

		baseURL, err := url.Parse(s.config.Service)
		if err != nil {
			return nil, err
		}

		meta, err := post(baseURL, s.config.Tags)
		if err != nil {
			return nil, err
		}

		meta.FilteredTags = intersect.SimpleGeneric(s.config.ExtractTags, meta.Tags)

		s.metadata = meta
	}

	return s.metadata, nil
}

func makeTags(elems []string) string {
	tagstring := ""
	for _, t := range elems {
		tagstring += "#" + strcase.ToCamel(t) + " "
	}
	return tagstring
}

func (s *Source) Caption() (string, error) {
	if s.config.Content == "" {
		return "", nil
	}

	md, err := s.getMeta()
	if err != nil {
		return "", err
	}

	t := template.New("post").Funcs(template.FuncMap{"MakeTags": makeTags})
	t, err = t.Parse(s.config.Content)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, md); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (s *Source) GetImageReader() (io.ReadCloser, string, error) {
	md, err := s.getMeta()
	if err != nil {
		return nil, "", err
	}

	res, err := get(md.Image)
	if err != nil {
		return nil, "", err
	}

	if res.StatusCode != 200 {
		return nil, "", ErrImageNotFound
	}

	return res.Body, path.Base(md.Image), nil
}

func (s *Source) GetMd5Hash() string {
	md, err := s.getMeta()
	if err != nil {
		return ""
	}
	return md.Md5
}

func (s *Source) IsSensitive() bool {
	return s.GetRating() != config.Safe
}

func (s *Source) GetRating() config.Rating {
	md, err := s.getMeta()
	if err != nil {
		return config.Safe
	}

	switch md.Rating {
	case "s":
	case "general":
		return config.Safe
	case "sensitive":
		return config.Sensitive
	case "q":
	case "questionable":
		return config.Questionable
	case "e":
	case "explicit":
		return config.Explicit
	}
	return config.Explicit
}

func (s *Source) GetTags() []string {
	md, err := s.getMeta()
	if err != nil {
		return []string{}
	}
	return md.Tags
}
