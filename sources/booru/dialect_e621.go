package booru

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/averagesecurityguy/random"
)

type e621Response struct {
	Posts []struct {
		// ApproverID   int64   `json:"approver_id"`
		// ChangeSeq    int64   `json:"change_seq"`
		// CommentCount int64   `json:"comment_count"`
		// CreatedAt    string  `json:"created_at"`
		// Description  string  `json:"description"`
		// Duration     float64 `json:"duration"`
		// FavCount     int64   `json:"fav_count"`
		File struct {
			// Ext    string `json:"ext"`
			// Height int64  `json:"height"`
			Md5 string `json:"md5"`
			// Size   int64  `json:"size"`
			URL string `json:"url"`
			// Width  int64  `json:"width"`
		} `json:"file"`
		//Flags struct {
		//	Deleted      bool `json:"deleted"`
		//	Flagged      bool `json:"flagged"`
		//	NoteLocked   bool `json:"note_locked"`
		//	Pending      bool `json:"pending"`
		//	RatingLocked bool `json:"rating_locked"`
		//	StatusLocked bool `json:"status_locked"`
		//} `json:"flags"`
		//HasNotes    bool          `json:"has_notes"`
		//ID          int64         `json:"id"`
		//IsFavorited bool          `json:"is_favorited"`
		//LockedTags  []interface{} `json:"locked_tags"`
		//Pools       []interface{} `json:"pools"`
		//Preview     struct {
		//	Height int64  `json:"height"`
		//	URL    string `json:"url"`
		//	Width  int64  `json:"width"`
		//} `json:"preview"`
		Rating string `json:"rating"`
		//Relationships struct {
		//	Children          []interface{} `json:"children"`
		//	HasActiveChildren bool          `json:"has_active_children"`
		//	HasChildren       bool          `json:"has_children"`
		//	ParentID          interface{}   `json:"parent_id"`
		//} `json:"relationships"`
		//Sample struct {
		//	Alternates struct {
		//		Original struct {
		//			Height int64         `json:"height"`
		//			Type   string        `json:"type"`
		//			Urls   []interface{} `json:"urls"`
		//			Width  int64         `json:"width"`
		//		} `json:"original"`
		//	} `json:"alternates"`
		//	Has    bool   `json:"has"`
		//	Height int64  `json:"height"`
		//	URL    string `json:"url"`
		//	Width  int64  `json:"width"`
		//} `json:"sample"`
		//Score struct {
		//	Down  int64 `json:"down"`
		//	Total int64 `json:"total"`
		//	Up    int64 `json:"up"`
		//} `json:"score"`
		Sources []string            `json:"sources"`
		Tags    map[string][]string `json:"tags"`
		// UpdatedAt  string `json:"updated_at"`
		// UploaderID int64  `json:"uploader_id"`
	} `json:"posts"`
}

func e621RandPage() uint64 {
	v, err := random.Uint64Range(0, 10)
	if err != nil {
		return 0
	}

	return v
}

func e621RandPost(max int) int {
	v, err := random.Uint64Range(0, uint64(max))
	if err != nil {
		return 0
	}

	return int(v)
}

func e621Post(base *url.URL, tags string) (*booruMetadata, error) {
	postURL, err := url.Parse(fmt.Sprintf("/posts.json?limit=100&&page=%v&tags=%v",
		e621RandPage(),
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

	parsed := new(e621Response)
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, ErrRequestFailed
	}

	post := parsed.Posts[gelbooruRandPost(len(parsed.Posts))]

	tagstring := ""
	for _, t := range post.Tags {
		for _, v := range t {
			tagstring += "#" + v + " "
		}
	}

	return &booruMetadata{
		Source:    post.Sources[0],
		Rating:    post.Rating,
		TagString: tagstring,
		Image:     post.File.URL,
	}, nil
}
