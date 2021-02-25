package booru

import (
	"encoding/json"
	"net/url"
	"strings"
)

type danbooruResponse struct {
	// ID           int    `json:"id"`
	Source       string `json:"source"`
	Rating       string `json:"rating"`
	TagString    string `json:"tag_string"`
	LargeFileURL string `json:"large_file_url"`
	// CreatedAt  string `json:"created_at"`
	// UploaderID int    `json:"uploader_id"`
	// Score      int    `json:"score"`
	// Md5        string `json:"md5"`
	// LastCommentBumpedAt interface{} `json:"last_comment_bumped_at"`
	// ImageWidth   int    `json:"image_width"`
	// ImageHeight  int    `json:"image_height"`
	// IsNoteLocked bool   `json:"is_note_locked"`
	// FavCount     int    `json:"fav_count"`
	// FileExt      string `json:"file_ext"`
	// LastNotedAt         interface{} `json:"last_noted_at"`
	// IsRatingLocked bool `json:"is_rating_locked"`
	// ParentID            interface{} `json:"parent_id"`
	// HasChildren        bool        `json:"has_children"`
	// ApproverID         int         `json:"approver_id"`
	// TagCountGeneral    int         `json:"tag_count_general"`
	// TagCountArtist     int         `json:"tag_count_artist"`
	// TagCountCharacter  int         `json:"tag_count_character"`
	// TagCountCopyright  int         `json:"tag_count_copyright"`
	// FileSize           int         `json:"file_size"`
	// IsStatusLocked     bool        `json:"is_status_locked"`
	// PoolString         string      `json:"pool_string"`
	// UpScore            int         `json:"up_score"`
	// DownScore          int         `json:"down_score"`
	// IsPending          bool        `json:"is_pending"`
	// IsFlagged          bool        `json:"is_flagged"`
	// IsDeleted          bool        `json:"is_deleted"`
	// TagCount           int         `json:"tag_count"`
	// UpdatedAt          string      `json:"updated_at"`
	// IsBanned           bool        `json:"is_banned"`
	// PixivID            int         `json:"pixiv_id"`
	// LastCommentedAt    interface{} `json:"last_commented_at"`
	// HasActiveChildren  bool        `json:"has_active_children"`
	// BitFlags           int         `json:"bit_flags"`
	// TagCountMeta       int         `json:"tag_count_meta"`
	// HasLarge           bool        `json:"has_large"`
	// HasVisibleChildren bool        `json:"has_visible_children"`
	// TagStringGeneral   string      `json:"tag_string_general"`
	// TagStringCharacter string      `json:"tag_string_character"`
	// TagStringCopyright string      `json:"tag_string_copyright"`
	// TagStringArtist    string      `json:"tag_string_artist"`
	// TagStringMeta      string      `json:"tag_string_meta"`
	// FileURL            string      `json:"file_url"`
	// PreviewFileURL     string      `json:"preview_file_url"`
}

func danbooruPost(base *url.URL, tags string) (*booruMetadata, error) {
	postURL, err := url.Parse("/posts/random.json?tags=" + url.QueryEscape(tags))
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

	parsed := new(danbooruResponse)
	if err := json.NewDecoder(res.Body).Decode(parsed); err != nil {
		return nil, ErrRequestFailed
	}

	tagArray := strings.Split(parsed.TagString, " ")
	tagstring := ""
	for _, t := range tagArray {
		tagstring += "#" + t + " "
	}

	return &booruMetadata{
		Source:    parsed.Source,
		Rating:    parsed.Rating,
		TagString: tagstring,
		Image:     parsed.LargeFileURL,
	}, nil
}
