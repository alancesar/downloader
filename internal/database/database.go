package database

import (
	"context"
	"errors"
	"github.com/alancesar/downloader/internal/database/internal"
	"gorm.io/gorm"
)

type (
	Gorm struct {
		db *gorm.DB
	}
)

func NewGorm(db *gorm.DB) (*Gorm, error) {
	if err := db.AutoMigrate(&internal.Download{}); err != nil {
		return nil, err
	}

	return &Gorm{
		db: db,
	}, nil
}

func (g Gorm) ExistByURL(ctx context.Context, url string) (bool, error) {
	download := internal.Download{}
	tx := g.db.WithContext(ctx).
		Where("url = ?", url).
		First(&download)

	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g Gorm) ExistsByETag(ctx context.Context, etag string) (bool, error) {
	download := internal.Download{}
	tx := g.db.WithContext(ctx).
		Where("e_tag = ?", etag).
		First(&download)

	if err := tx.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g Gorm) Save(ctx context.Context, url, eTag string) error {
	return g.db.WithContext(ctx).Save(&internal.Download{
		ETag: eTag,
		URL:  url,
	}).Error
}
