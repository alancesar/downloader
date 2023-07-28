package redgifs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/downloader/pkg/redgifs/param"
	"github.com/alancesar/downloader/pkg/redgifs/testdata"
	"github.com/alancesar/downloader/pkg/transport"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_GetGIFByID(t *testing.T) {
	type fields struct {
		HTTPClient *http.Client
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    GIF
		wantErr bool
	}{
		{
			name: "Should retrieve GIF by ID properly",
			fields: fields{
				HTTPClient: testdata.NewHTTPClient([]byte(testdata.RedGIFsGifResponse), http.StatusOK, nil),
			},
			args: args{
				id: "some-media",
			},
			want: GIF{
				ID:          "some-media",
				Width:       1920,
				Height:      1080,
				Niches:      []string{"some-niche"},
				Tags:        []string{"Some Tag"},
				Verified:    true,
				Description: "Some description",
				Published:   true,
				URLs: Source{
					Thumbnail:  "https://thumbs44.redgifs.com/some-media-mobile.jpg?expires=1680625800&signature=v2:755a40e0b354717fc9cfa250cc01695e89ce2b830d9f30973fd5498acace10a1&for=192.198.0.1&hash=6163438793",
					VThumbnail: "https://thumbs44.redgifs.com/some-media-mobile.mp4?expires=1680625800&signature=v2:a7c3be8860e39bba9e69ebc96f3bd288cd9e251647a333ea18e86e28b52c7b95&for=192.198.0.1&hash=6163438793",
					Poster:     "https://thumbs44.redgifs.com/some-media-poster.jpg?expires=1680625800&signature=v2:761f3ab07837f9066aa605237fd3f82fa9f0097b65490ca5af892f40a84b1016&for=192.198.0.1&hash=6163438793",
					SD:         "https://thumbs44.redgifs.com/some-media-mobile.mp4?expires=1680625800&signature=v2:a7c3be8860e39bba9e69ebc96f3bd288cd9e251647a333ea18e86e28b52c7b95&for=192.198.0.1&hash=6163438793",
					HD:         "https://thumbs44.redgifs.com/some-media.mp4?expires=1680625800&signature=v2:d330e63fd8fbf4dda3f7d66d399b9c1fffdc7383c1190bee65f694b9f5affb29&for=192.198.0.1&hash=6163438793",
				},
				UserName:  "some-username",
				Gallery:   "some-gallery",
				Sexuality: []string{"straight"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.HTTPClient)
			got, err := c.GetGIFByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGIFByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGIFByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RetrieveToken(t *testing.T) {
	type fields struct {
		httpClient *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should retrieve token properly",
			fields: fields{
				httpClient: testdata.NewHTTPClient([]byte(testdata.RedGIFsTokenResponse), http.StatusOK, nil),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    "some.token.here",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.httpClient)
			got, err := c.RetrieveToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RetrieveToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_AuthorSearchStream(t *testing.T) {
	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	redGIFsAuthProvider := NewClient(client)
	redGIFsAuthClient := &http.Client{
		Transport: transport.NewAuthorizationRoundTripper(func() string {
			token, err := redGIFsAuthProvider.RetrieveToken()
			if err != nil {
				return ""
			}

			return fmt.Sprintf("Bearer %s", token)
		}, client.Transport),
	}
	redGIFsClient := NewClient(redGIFsAuthClient)
	stream := redGIFsClient.AuthorSearchStream(context.Background(), param.AuthorSearch{
		Username: "alexasweetie00",
	})

	for searchResult := range stream {
		for _, gif := range searchResult.GIFs {
			body := new(bytes.Buffer)
			gifURL := fmt.Sprintf("https://www.redgifs.com/watch/%s", gif.ID)
			_ = json.NewEncoder(body).Encode(struct {
				URL string `json:"url"`
			}{
				URL: gifURL,
			})
			res, _ := http.Post("https://redgifs-fetcher-347bebb7b806.herokuapp.com/publish", "application/json", body)
			if res.StatusCode != http.StatusAccepted {
				log.Println("bad status", res.StatusCode, gifURL)
			}
		}
	}
}
