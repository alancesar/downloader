package usecase

import (
	"context"
	"errors"
	"github.com/alancesar/downloader/pkg/media"
	"github.com/alancesar/downloader/pkg/status"
)

type (
	Database interface {
		ExistByURL(ctx context.Context, url string) (bool, error)
		ExistsByETag(ctx context.Context, etag string) (bool, error)
		Save(ctx context.Context, url, etag string) error
	}

	Downloader interface {
		Download(ctx context.Context, m media.Media) error
		GetMetadata(m media.Media) (media.Metadata, error)
	}

	Interceptor interface {
		Intercept(ctx context.Context, m media.Media) (media.Media, error)
	}

	Download struct {
		db           Database
		downloader   Downloader
		interceptors map[string]Interceptor
	}
)

func NewDownload(downloader Downloader, db Database, interceptors map[string]Interceptor) *Download {
	return &Download{
		downloader:   downloader,
		db:           db,
		interceptors: interceptors,
	}
}

func (d Download) Execute(ctx context.Context, m media.Media, provider string) error {
	if exists, err := d.db.ExistByURL(ctx, m.URL); err != nil {
		return err
	} else if exists {
		return nil
	}

	if provider != "" {
		var err error
		if interceptor, ok := d.interceptors[provider]; ok {
			if m, err = interceptor.Intercept(ctx, m); err != nil {
				if errors.Is(err, status.ErrNotFound) {
					return nil
				}

				return err
			}
		}
	}

	md, err := d.downloader.GetMetadata(m)
	if err != nil {
		if errors.Is(err, status.ErrNotFound) {
			return d.db.Save(ctx, m.URL, "")
		}
		return err
	}

	if md.ETag != "" {
		if exists, err := d.db.ExistsByETag(ctx, md.ETag); err != nil {
			return err
		} else if exists {
			return d.db.Save(ctx, m.URL, md.ETag)
		}
	}

	if m.Filename == "" {
		m.Filename = m.BuildFilename(md.ContentType)
	}

	if err := d.downloader.Download(ctx, m); err != nil {
		return err
	}

	return d.db.Save(ctx, m.URL, md.ETag)
}
