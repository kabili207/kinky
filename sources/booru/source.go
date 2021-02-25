package booru

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
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

		s.metadata = meta
	}

	return s.metadata, nil
}

func (s *Source) Caption() (string, error) {
	if s.config.Content == "" {
		return "", nil
	}

	md, err := s.getMeta()
	if err != nil {
		return "", err
	}

	t := template.New("post")
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

func (s *Source) GetImageReader() (io.ReadCloser, error) {
	md, err := s.getMeta()
	if err != nil {
		return nil, err
	}

	res, err := get(md.Image)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrImageNotFound
	}

	return res.Body, nil
}

func (s *Source) IsSensitive() bool {
	md, err := s.getMeta()
	if err != nil {
		return false
	}

	return md.Rating != "s"
}
