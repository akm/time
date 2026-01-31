package faketime

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileProvider(t *testing.T) {
	filePath := "/some/path/to/file"
	provider := NewFileProvider(filePath)

	if provider.filePath != filePath {
		t.Errorf("filePath = %v, want %v", provider.filePath, filePath)
	}
}

func TestFileProvider_Get(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, dir string) string
		want        string
		wantErr     error
		wantErrMsg  string
	}{
		{
			name: "file does not exist returns empty string",
			setup: func(t *testing.T, dir string) string {
				return filepath.Join(dir, "nonexistent.txt")
			},
			want:    "",
			wantErr: nil,
		},
		{
			name: "path is a directory returns error",
			setup: func(t *testing.T, dir string) string {
				subdir := filepath.Join(dir, "subdir")
				if err := os.Mkdir(subdir, 0755); err != nil {
					t.Fatal(err)
				}
				return subdir
			},
			wantErr: ErrFileRead,
		},
		{
			name: "file with content",
			setup: func(t *testing.T, dir string) string {
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte("2024-01-02 15:04:05"), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			want: "2024-01-02 15:04:05",
		},
		{
			name: "file with whitespace is trimmed",
			setup: func(t *testing.T, dir string) string {
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte("  2024-01-02 15:04:05  \n"), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			want: "2024-01-02 15:04:05",
		},
		{
			name: "file with ratio option",
			setup: func(t *testing.T, dir string) string {
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte("2024-01-02 15:04:05 x2\n"), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			want: "2024-01-02 15:04:05 x2",
		},
		{
			name: "empty file",
			setup: func(t *testing.T, dir string) string {
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			want: "",
		},
		{
			name: "file with only whitespace",
			setup: func(t *testing.T, dir string) string {
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte("   \n\t  "), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			filePath := tt.setup(t, dir)

			provider := NewFileProvider(filePath)
			got, err := provider.Get(context.Background())

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Get() = %q, want %q", got, tt.want)
			}
		})
	}
}
