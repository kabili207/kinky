package booru

import (
	"context"
	"net/http"
	"net/url"
)

type dialectPost func(base *url.URL, tags string) (*booruMetadata, error)

var dialects map[string]dialectPost = map[string]dialectPost{
	"danbooru": danbooruPost, // danbooru, safebooru
	"gelbooru": gelbooruPost, // gelbooru, rule34.xxx, xbooru.com
	"e621":     e621Post,     // e621.net
	"moebooru": moebooruPost, // konachan.net, yande.re
}

type booruMetadata struct {
	Source    string
	Rating    string
	TagString string
	Image     string
}

func get(url string) (*http.Response, error) {
	r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(r)
}
