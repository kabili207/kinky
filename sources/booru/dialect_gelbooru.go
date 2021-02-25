package booru

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/averagesecurityguy/random"
)

type gelbooruResponse struct {
	Source    string `json:"source"`
	Tags      string `json:"tags"`
	FileURL   string `json:"file_url"`
	Rating    string `json:"rating"`
	Directory string `json:"directory"`
	Image     string `json:"image"`

	// Change        int64       `json:"change"`
	// CreatedAt     string      `json:"created_at"`
	// Hash          string      `json:"hash"`
	// Height        int64       `json:"height"`
	// ID            int64       `json:"id"`
	// Owner         string      `json:"owner"`
	// ParentID      interface{} `json:"parent_id"`
	// PreviewHeight int64       `json:"preview_height"`
	// PreviewWidth  int64       `json:"preview_width"`
	// Sample        int64       `json:"sample"`
	// SampleHeight  int64       `json:"sample_height"`
	// SampleWidth   int64       `json:"sample_width"`
	// Score         int64       `json:"score"`
	// Title         string      `json:"title"`
	// Width         int64       `json:"width"`
}

func gelbooruRandPage() uint64 {
	v, err := random.Uint64Range(0, 10)
	if err != nil {
		return 0
	}

	return v
}

func gelbooruRandPost(max int) int {
	v, err := random.Uint64Range(0, uint64(max))
	if err != nil {
		return 0
	}

	return int(v)
}

func gelbooruPost(base *url.URL, tags string) (*booruMetadata, error) {
	postURL, err := url.Parse(fmt.Sprintf("/index.php?page=dapi&s=post&q=index&json=1&limit=100&&pid=%v&tags=%v",
		gelbooruRandPage(),
		url.QueryEscape(tags)))
	if err != nil {
		return nil, err
	}
	postURL = base.ResolveReference(postURL)
	res, err := get(postURL.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, ErrRequestFailed
	}
	defer res.Body.Close()

	var parsed []gelbooruResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, ErrRequestFailed
	}

	post := parsed[gelbooruRandPost(len(parsed))]

	tagArray := strings.Split(post.Tags, " ")
	tagstring := ""
	for _, t := range tagArray {
		tagstring += "#" + t + " "
	}

	image := post.FileURL

	if image == "" {
		image = fmt.Sprintf("%v/images/%v/%v",
			base,
			post.Directory,
			post.Image)
	}

	return &booruMetadata{
		Source:    post.Source,
		Rating:    post.Rating,
		TagString: tagstring,
		Image:     image,
	}, nil
}