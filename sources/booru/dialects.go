package booru

import (
	"context"
	"net/http"
	"net/url"
)

type dialect interface {
	randomPost(base *url.URL, tags string) (*booruMetadata, error)
}

var dialects map[string]dialect = map[string]dialect{
	"danbooru": newDanbooruDialect(),
}

type booruMetadata struct {
	Source       string `json:"source"`
	Rating       string `json:"rating"`
	TagString    string `json:"tag_string"`
	LargeFileURL string `json:"large_file_url"`
}

func get(url string) (*http.Response, error) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}
