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
		ExistByETag(ctx context.Context, etag string) (bool, error)
		Save(ctx context.Context, url, etag string) error
	}

	Downloader interface {
		Download(ctx context.Context, m media.Media) error
		GetMetadata(m media.Media) (media.Metadata, error)
	}

	Download struct {
		db         Database
		downloader Downloader
	}
)

func NewDownload(downloader Downloader, db Database) *Download {
	return &Download{
		downloader: downloader,
		db:         db,
	}
}

func (d Download) Execute(ctx context.Context, m media.Media) error {
	cleanURL := m.CleanURL()
	if exists, err := d.db.ExistByURL(ctx, cleanURL); err != nil {
		return err
	} else if exists {
		return nil
	}

	md, err := d.downloader.GetMetadata(m)
	if err != nil {
		if errors.Is(err, status.ErrNotFound) {
			return d.db.Save(ctx, cleanURL, "")
		}
		return err
	}

	if md.ETag != "" {
		if exists, err := d.db.ExistByETag(ctx, md.ETag); err != nil {
			return err
		} else if exists {
			return d.db.Save(ctx, cleanURL, md.ETag)
		}
	}

	if m.Filename == "" {
		m.Filename = m.BuildFilename(md.ContentType)
	}

	if err := d.downloader.Download(ctx, m); err != nil {
		return err
	}

	return d.db.Save(ctx, cleanURL, md.ETag)
}
