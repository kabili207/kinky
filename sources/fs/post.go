package fs

import "net/url"

type Post struct {
	nsfw    *bool
	image   *url.URL
	content *string
}

func NewPost(url *url.URL) *Post {
	return &Post{
		image: url,
	}
}

func (p *Post) SetNSFW(flag bool) *Post {
	p.nsfw = &flag

	return p
}

func (p *Post) ClearNSFW() *Post {
	p.nsfw = nil

	return p
}

func (p *Post) SetContent(content string) *Post {
	p.content = &content

	return p
}

func (p *Post) ClearContent() *Post {
	p.content = nil

	return p
}
