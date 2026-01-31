package faketime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akm/time"
)

func TestNewFile(t *testing.T) {
	filePath := "/tmp/test.txt"
	layout := "2006-01-02 15:04:05"

	f := NewFile(filePath, layout)

	if f.FilePath != filePath {
		t.Errorf("FilePath = %v, want %v", f.FilePath, filePath)
	}
	if f.layout != layout {
		t.Errorf("layout = %v, want %v", f.layout, layout)
	}
}

func TestFile_Save(t *testing.T) {
	tests := []struct {
		name        string
		layout      string
		time        time.Time
		wantContent string
	}{
		{
			name:        "standard layout",
			layout:      "2006-01-02 15:04:05",
			time:        time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC),
			wantContent: "2024-06-15 12:30:45",
		},
		{
			name:        "date only",
			layout:      "2006-01-02",
			time:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			wantContent: "2024-01-02",
		},
		{
			name:        "RFC3339",
			layout:      time.RFC3339,
			time:        time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC),
			wantContent: "2024-06-15T12:30:45Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "faketime.txt")

			f := NewFile(filePath, tt.layout)
			if err := f.Save(tt.time); err != nil {
				t.Fatalf("Save() error = %v", err)
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.wantContent {
				t.Errorf("file content = %v, want %v", string(content), tt.wantContent)
			}
		})
	}
}

func TestFile_Save_Overwrite(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "faketime.txt")
	layout := "2006-01-02 15:04:05"

	f := NewFile(filePath, layout)

	// Save first time
	time1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := f.Save(time1); err != nil {
		t.Fatalf("Save() first call error = %v", err)
	}

	// Save second time (should overwrite)
	time2 := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	if err := f.Save(time2); err != nil {
		t.Fatalf("Save() second call error = %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	wantContent := "2024-12-31 23:59:59"
	if string(content) != wantContent {
		t.Errorf("file content = %v, want %v", string(content), wantContent)
	}
}

func TestFile_Save_InvalidPath(t *testing.T) {
	f := NewFile("/nonexistent/directory/faketime.txt", "2006-01-02 15:04:05")
	err := f.Save(time.Now())
	if err == nil {
		t.Error("Save() expected error for invalid path, got nil")
	}
}

func TestFile_Delete(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
	}{
		{
			name:       "delete existing file",
			fileExists: true,
		},
		{
			name:       "delete non-existing file",
			fileExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "faketime.txt")

			f := NewFile(filePath, "2006-01-02 15:04:05")

			if tt.fileExists {
				if err := f.Save(time.Now()); err != nil {
					t.Fatalf("Save() error = %v", err)
				}
				// Verify file exists
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Fatal("file should exist before delete")
				}
			}

			if err := f.Delete(); err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			// Verify file does not exist
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				t.Error("file should not exist after delete")
			}
		})
	}
}
