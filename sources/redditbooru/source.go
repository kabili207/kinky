package redditbooru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/averagesecurityguy/random"
)

type redditPost struct {
	CdnURL     string `json:"cdnUrl"`
	Title      string `json:"title"`
	ExternalID string `json:"externalId"`
	Nsfw       bool   `json:"nsfw"`
	// Caption    interface{} `json:"caption"`
	// SourceURL  interface{} `json:"sourceUrl"`
	// ID          int         `json:"id"`
	// ImageID     int         `json:"imageId"`
	// Width       int         `json:"width"`
	// Height      int         `json:"height"`
	// Type        string      `json:"type"`
	// SourceID    int         `json:"sourceId"`
	// SourceName  string      `json:"sourceName"`
	// PostID      int         `json:"postId"`
	// Keywords    string      `json:"keywords"`
	// DateCreated int         `json:"dateCreated"`
	// Score       string      `json:"score"`
	// Visible     bool        `json:"visible"`
	// UserID      int         `json:"userId"`
	// UserName    string      `json:"userName"`
	// Thumb       string      `json:"thumb"`
	// IdxInAlbum  interface{} `json:"idxInAlbum"`
	// Age         int         `json:"age"`
}

type Source struct {
	config *SourceConfig
	post   *redditPost
}

func get(url string) (*http.Response, error) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}

func New(cfg *SourceConfig) (*Source, error) {
	s := new(Source)
	s.config = cfg

	return s, nil
}

func (s *Source) getPost() (*redditPost, error) {
	if s.post == nil {
		sources := ""
		for _, filter := range s.config.Filter {
			sources += fmt.Sprintf("%v,", filter)
		}

		res, err := get(fmt.Sprintf("https://redditbooru.com/images/?sources=%v", url.QueryEscape(sources)))
		if err != nil {
			return nil, err
		}

		if res.StatusCode != 200 {
			return nil, ErrRequestFailed
		}
		defer res.Body.Close()

		var parsed []redditPost
		if err = json.NewDecoder(res.Body).Decode(&parsed); err != nil {
			return nil, ErrRequestFailed
		}

		index, err := random.Uint64Range(0, uint64(len(parsed)-1))
		if err != nil {
			return nil, err
		}

		s.post = &parsed[index]
	}

	return s.post, nil
}

func (s *Source) Caption() (string, error) {
	post, err := s.getPost()
	if err != nil {
		return "", err
	}

	return post.Title, nil
}

func (s *Source) GetImageReader() (io.ReadCloser, error) {
	post, err := s.getPost()
	if err != nil {
		return nil, err
	}

	res, err := get(post.CdnURL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrImageNotFound
	}

	return res.Body, nil
}

func (s *Source) IsSensitive() bool {
	post, err := s.getPost()
	if err != nil {
		return false
	}

	return post.Nsfw
}
