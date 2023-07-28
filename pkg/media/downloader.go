package media

import (
	"context"
	"fmt"
	"github.com/alancesar/downloader/pkg/status"
	"io"
	"net/http"
)

const (
	suffix = ".tmp"
)

type (
	Storage interface {
		Create(ctx context.Context, path string) (io.WriteCloser, error)
		Exist(ctx context.Context, path string) (bool, error)
		Rename(ctx context.Context, src, dst string) error
		Remove(ctx context.Context, path string) error
	}

	ProgressBar func(response *http.Response, filename string) io.Writer

	Downloader struct {
		storage Storage
		pb      ProgressBar
		client  *http.Client
	}
)

func NewDownloader(storage Storage, pb ProgressBar, client *http.Client) *Downloader {
	return &Downloader{
		storage: storage,
		pb:      pb,
		client:  client,
	}
}

func (d Downloader) Download(ctx context.Context, media Media) error {
	if exist, err := d.storage.Exist(ctx, media.Path()); err != nil {
		return err
	} else if exist {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, media.URL, nil)
	if err != nil {
		return err
	}

	res, err := d.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w: %s", status.ErrNotFound, media.URL)
	} else if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %d (%s): %s", status.ErrBadStatus, res.StatusCode, res.Status, media.URL)
	}

	return d.write(ctx, res, media)
}

func (d Downloader) write(ctx context.Context, response *http.Response, media Media) error {
	temp := media.Path() + suffix
	writer, err := d.storage.Create(ctx, temp)
	if err != nil {
		return err
	}

	done := make(chan error)

	go func() {
		multiWriter := io.MultiWriter(writer, d.pb(response, fmt.Sprintf("%32.32s", media.Filename)))
		if _, err = io.Copy(multiWriter, response.Body); err != nil {
			_ = d.storage.Remove(ctx, temp)
			done <- err
			return
		}

		done <- d.storage.Rename(ctx, temp, media.Path())
	}()

	for {
		select {
		case <-ctx.Done():
			_ = writer.Close()
			_ = d.storage.Remove(ctx, temp)
			return ctx.Err()
		case err := <-done:
			_ = writer.Close()
			return err
		}
	}
}

func (d Downloader) GetMetadata(media Media) (Metadata, error) {
	res, err := d.client.Head(media.URL)
	if err != nil {
		return Metadata{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return Metadata{}, fmt.Errorf("%w: %s", status.ErrNotFound, media.URL)
	} else if res.StatusCode >= http.StatusBadRequest {
		return Metadata{}, fmt.Errorf("%w: %d (%s): %s", status.ErrBadStatus, res.StatusCode, res.Status, media.URL)
	}

	return Metadata{
		ContentType: ContentType(res.Header.Get("Content-Type")),
		ETag:        res.Header.Get("ETag"),
	}, nil
}
