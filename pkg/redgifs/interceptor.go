package redgifs

import (
	"context"
	"github.com/alancesar/downloader/pkg/media"
)

type (
	Interceptor struct {
		client *Client
	}
)

func NewInterceptor(client *Client) *Interceptor {
	return &Interceptor{
		client: client,
	}
}

func (i Interceptor) Intercept(_ context.Context, m media.Media) (media.Media, error) {
	gif, err := i.client.GetGIFByURL(m.URL)
	if err != nil {
		return media.Media{}, err
	}

	return gif.ToMedia(), nil
}
