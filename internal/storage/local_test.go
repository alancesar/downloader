package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocal_Exist(t *testing.T) {
	root := t.TempDir()

	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		before  func(path string) error
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should return true if the file already exists",
			before: func(path string) error {
				_ = os.Remove(path)
				_, err := os.Create(path)
				return err
			},
			args: args{
				ctx:  context.TODO(),
				path: "some-file.txt",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should return false if the file do not exist",
			before: func(path string) error {
				_ = os.Remove(path)
				return nil
			},
			args: args{
				ctx:  context.TODO(),
				path: "some-file.txt",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completePath := buildCompletePath(root, tt.args.path)
			if err := tt.before(completePath); err != nil {
				t.Error(err)
				return
			}

			l := NewLocalStorage(root)
			got, err := l.Exist(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPublished() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsPublished() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocal_Create(t *testing.T) {
	root := t.TempDir()

	type args struct {
		ctx      context.Context
		filename string
		parent   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should create file at root path from storage",
			args: args{
				ctx:      context.TODO(),
				filename: "some-file.txt",
			},
			wantErr: false,
		},
		{
			name: "Should create file inside a directory of storage",
			args: args{
				ctx:      context.TODO(),
				filename: "some-file.txt",
				parent:   []string{"a", "b"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalStorage(root)
			_, err := l.Create(tt.args.ctx, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			completePath := buildCompletePath(root, tt.args.filename)
			_, err = os.Stat(completePath)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestLocal_Remove(t *testing.T) {
	root := t.TempDir()

	type args struct {
		ctx      context.Context
		filename string
		parent   []string
	}
	tests := []struct {
		name    string
		before  func(path string) error
		args    args
		wantErr bool
	}{
		{
			name: "Should remove file at root path from storage",
			before: func(path string) error {
				_, err := os.Create(path)
				return err
			},
			args: args{
				ctx:      context.TODO(),
				filename: "some-file.txt",
			},
			wantErr: false,
		},
		{
			name: "Should remove file inside a directory of storage",
			before: func(path string) error {
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					return err
				}
				_, err := os.Create(path)
				return err
			},
			args: args{
				ctx:      context.TODO(),
				filename: "some-file.txt",
				parent:   []string{"a", "b"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completePath := buildCompletePath(root, tt.args.filename)
			if err := tt.before(completePath); err != nil {
				t.Error(err)
				return
			}

			l := NewLocalStorage(root)
			if err := l.Remove(tt.args.ctx, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func buildCompletePath(root, path string) string {
	return filepath.Join(root, path)
}
