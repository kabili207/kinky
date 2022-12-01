package booru

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/averagesecurityguy/random"
)

type gelbooruContainer struct {
	Post []gelbooruResponse `json:"post"`
}

type gelbooruResponse struct {
	Source    string `json:"source"`
	Tags      string `json:"tags"`
	FileURL   string `json:"file_url"`
	Rating    string `json:"rating"`
	Directory string `json:"directory"`
	Image     string `json:"image"`
	Md5       string `json:"md5"`
	Title     string `json:"title"`

	// Change        int64       `json:"change"`
	// CreatedAt     string      `json:"created_at"`
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
	postURL, err := url.Parse(fmt.Sprintf("/index.php?page=dapi&s=post&q=index&json=1&limit=100&&pid=%v&tags=%v+sort:random",
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

	var parsed gelbooruContainer
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, ErrRequestFailed
	}

	post := parsed.Post[gelbooruRandPost(len(parsed.Post))]

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
		Source:    html.UnescapeString(post.Source),
		Rating:    post.Rating,
		TagString: tagstring,
		Image:     image,
		Md5:       post.Md5,
		Title:     html.UnescapeString(post.Title),
	}, nil
}
