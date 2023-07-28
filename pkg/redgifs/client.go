package redgifs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/downloader/pkg/redgifs/param"
	"github.com/alancesar/downloader/pkg/status"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	apiBasePath = "https://api.redgifs.com/v2"
)

type (
	Client struct {
		httpClient *http.Client
	}

	AuthProvider interface {
		RetrieveToken(ctx context.Context) (string, error)
	}

	Response struct {
		GIF    GIF      `json:"gif"`
		User   string   `json:"user"`
		Niches []string `json:"niches"`
	}

	SearchResult struct {
		Page  int   `json:"page"`
		Pages int   `json:"pages"`
		Total int   `json:"total"`
		GIFs  []GIF `json:"gifs"`
	}

	Auth struct {
		Addr  string `json:"addr"`
		Agent string `json:"agent"`
		Rtfm  string `json:"rtfm"`
		Token string `json:"token"`
	}
)

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c Client) GetGIFByURL(rawURL string) (GIF, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return GIF{}, err
	}

	base := path.Base(parsedURL.Path)
	id := strings.ToLower(base)
	return c.GetGIFByID(id)
}

func (c Client) GetGIFByID(id string) (GIF, error) {
	var response Response
	rawURL := fmt.Sprintf("%s/gifs/%s", apiBasePath, id)
	err := c.doGet(rawURL, &response)
	return response.GIF, err
}

func (c Client) AuthorSearch(params param.Search) (SearchResult, error) {
	var search SearchResult
	searchURL := apiBasePath + params.ToURL()
	err := c.doGet(searchURL, &search)
	return search, err
}

func (c Client) AuthorSearchStream(ctx context.Context, params param.AuthorSearch) <-chan SearchResult {
	stream := make(chan SearchResult)

	go func() {
		c.authorSearchStream(ctx, params, stream)
		close(stream)
	}()
	return stream
}

func (c Client) authorSearchStream(ctx context.Context, params param.AuthorSearch, stream chan<- SearchResult) {
	if ctx.Err() != nil {
		return
	}

	search, err := c.AuthorSearch(params)
	if err != nil {
		log.Println(err)
		return
	}

	stream <- search

	if search.Pages > params.Page {
		params.Page = params.Page + 1
		c.authorSearchStream(ctx, params, stream)
	}
}

func (c Client) RetrieveToken() (string, error) {
	var auth Auth
	authURL := fmt.Sprintf("%s/auth/temporary", apiBasePath)
	err := c.doGet(authURL, &auth)
	return auth.Token, err
}

func (c Client) doGet(url string, output any) error {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == http.StatusNotFound ||
		res.StatusCode == http.StatusGone {
		return fmt.Errorf("%w: %s", status.ErrNotFound, url)
	} else if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %d (%s): %s", status.ErrBadStatus, res.StatusCode, res.Status, url)
	}

	return json.NewDecoder(res.Body).Decode(&output)
}
