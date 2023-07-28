package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type (
	Local struct {
		root string
	}
)

func NewLocalStorage(root string) *Local {
	return &Local{
		root: root,
	}
}

func (l Local) Create(_ context.Context, path string) (io.WriteCloser, error) {
	completePath := l.buildCompletePath(path)
	dir := filepath.Dir(completePath)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return nil, err
	}

	return os.Create(completePath)
}

func (l Local) Remove(_ context.Context, path string) error {
	completePath := l.buildCompletePath(path)
	return os.Remove(completePath)
}

func (l Local) Exist(_ context.Context, path string) (bool, error) {
	completePath := l.buildCompletePath(path)
	if _, err := os.Stat(completePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("error on stat %s: %w", completePath, err)
	}

	return true, nil
}

func (l Local) buildCompletePath(path string) string {
	return filepath.Join(l.root, path)
}
