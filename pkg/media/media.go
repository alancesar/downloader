package media

import (
	"net/url"
	"path"
	"strings"
)

type (
	Provider    string
	ContentType string
	Status      string

	Media struct {
		URL      string   `json:"url"`
		Filename string   `json:"filename"`
		Parent   []string `json:"parent"`
	}

	Metadata struct {
		ContentType ContentType
		ETag        string
	}
)

func (t ContentType) Ext() string {
	_, ext, found := strings.Cut(string(t), "/")
	if found {
		return "." + ext
	}

	return ""
}

func (m Media) BuildFilename(contentType ContentType) string {
	parsedURL, _ := url.Parse(m.URL)
	filename := path.Base(parsedURL.Path)
	if ext := path.Ext(filename); ext != "" {
		return strings.ReplaceAll(filename, ext, contentType.Ext())
	}

	return filename + contentType.Ext()
}

func (m Media) Path() string {
	return path.Join(path.Join(m.Parent...), m.Filename)
}
