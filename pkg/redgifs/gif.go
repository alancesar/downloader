package redgifs

import (
	"github.com/alancesar/downloader/pkg/media"
	"net/url"
	"path"
)

type (
	GIF struct {
		ID          string   `json:"id"`
		Width       int      `json:"width"`
		Height      int      `json:"height"`
		Niches      []string `json:"niches"`
		Tags        []string `json:"tags"`
		Verified    bool     `json:"verified"`
		Description string   `json:"description"`
		Published   bool     `json:"published"`
		URLs        Source   `json:"urls"`
		UserName    string   `json:"userName"`
		Gallery     string   `json:"gallery"`
		Sexuality   []string `json:"sexuality"`
		CreateDate  float64  `json:"create_date"`
	}

	Source struct {
		Thumbnail  string `json:"thumbnail"`
		VThumbnail string `json:"vthumbnail"`
		Poster     string `json:"poster"`
		SD         string `json:"sd"`
		HD         string `json:"hd"`
	}
)

func (g GIF) ToMedia() media.Media {
	return media.Media{
		URL:      g.URLs.HD,
		Parent:   []string{"redgifs", g.UserName},
		Filename: g.Filename(),
	}
}

func (g GIF) Filename() string {
	parsedURL, _ := url.Parse(g.URLs.HD)
	return path.Base(parsedURL.Path)
}
