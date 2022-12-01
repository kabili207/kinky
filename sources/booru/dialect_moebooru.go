package booru

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/averagesecurityguy/random"
)

type moebooruResponse struct {
	FileURL string `json:"file_url"`
	Rating  string `json:"rating"`
	Tags    string `json:"tags"`
	Source  string `json:"source"`
	Md5     string `json:"md5"`

	// ActualPreviewHeight int64         `json:"actual_preview_height"`
	// ActualPreviewWidth  int64         `json:"actual_preview_width"`
	// Author              string        `json:"author"`
	// Change              int64         `json:"change"`
	// CreatedAt           int64         `json:"created_at"`
	// CreatorID           int64         `json:"creator_id"`
	// FileSize            int64         `json:"file_size"`
	// FlagDetail          interface{}   `json:"flag_detail"`
	// Frames              []interface{} `json:"frames"`
	// FramesPending       []interface{} `json:"frames_pending"`
	// FramesPendingString string        `json:"frames_pending_string"`
	// FramesString        string        `json:"frames_string"`
	// HasChildren         bool          `json:"has_children"`
	// Height              int64         `json:"height"`
	// ID                  int64         `json:"id"`
	// IsHeld              bool          `json:"is_held"`
	// IsShownInIndex      bool          `json:"is_shown_in_index"`
	// JpegFileSize        int64         `json:"jpeg_file_size"`
	// JpegHeight          int64         `json:"jpeg_height"`
	// JpegURL             string        `json:"jpeg_url"`
	// JpegWidth           int64         `json:"jpeg_width"`
	// ParentID            int64         `json:"parent_id"`
	// PreviewHeight       int64         `json:"preview_height"`
	// PreviewURL          string        `json:"preview_url"`
	// PreviewWidth        int64         `json:"preview_width"`
	// SampleFileSize      int64         `json:"sample_file_size"`
	// SampleHeight        int64         `json:"sample_height"`
	// SampleURL           string        `json:"sample_url"`
	// SampleWidth         int64         `json:"sample_width"`
	// Score               int64         `json:"score"`
	// Status              string        `json:"status"`
	// Width               int64         `json:"width"`
}

func moebooruRandPage() uint64 {
	v, err := random.Uint64Range(0, 10)
	if err != nil {
		return 0
	}

	return v
}

func moebooruRandPost(max int) int {
	v, err := random.Uint64Range(0, uint64(max))
	if err != nil {
		return 0
	}

	return int(v)
}

func moebooruPost(base *url.URL, tags string) (*booruMetadata, error) {
	postURL, err := url.Parse(fmt.Sprintf("/post.json?page=%v&tags=%v",
		moebooruRandPage(),
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

	var parsed []moebooruResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, ErrRequestFailed
	}

	post := parsed[moebooruRandPost(len(parsed))]

	tagArray := strings.Split(post.Tags, " ")
	tagstring := ""
	for _, t := range tagArray {
		tagstring += "#" + t + " "
	}

	return &booruMetadata{
		Source:    post.Source,
		Rating:    post.Rating,
		TagString: tagstring,
		Image:     post.FileURL,
	}, nil
}
