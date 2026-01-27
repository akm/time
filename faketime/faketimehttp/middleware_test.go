package faketimehttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/akm/time"
)

func TestMiddleware(t *testing.T) {
	t.Run("file conditions", func(t *testing.T) {
		tests := []struct {
			name      string
			setup     func(t *testing.T) string // returns file path
			wantCode  int
			checkTime bool // whether time should be modified
		}{
			{
				name: "file does not exist",
				setup: func(t *testing.T) string {
					return "/nonexistent/path/to/file"
				},
				wantCode:  http.StatusOK,
				checkTime: false,
			},
			{
				name: "path is a directory",
				setup: func(t *testing.T) string {
					return t.TempDir()
				},
				wantCode:  http.StatusOK,
				checkTime: false,
			},
			{
				name: "file is empty",
				setup: func(t *testing.T) string {
					dir := t.TempDir()
					filePath := filepath.Join(dir, "time.txt")
					if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
						t.Fatal(err)
					}
					return filePath
				},
				wantCode:  http.StatusOK,
				checkTime: false,
			},
			{
				name: "file contains only whitespace",
				setup: func(t *testing.T) string {
					dir := t.TempDir()
					filePath := filepath.Join(dir, "time.txt")
					if err := os.WriteFile(filePath, []byte("  \n\t  "), 0644); err != nil {
						t.Fatal(err)
					}
					return filePath
				},
				wantCode:  http.StatusOK,
				checkTime: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filePath := tt.setup(t)
				handler := Middleware(filePath, time.RFC3339)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
					}),
				)

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if rec.Code != tt.wantCode {
					t.Errorf("got status %d, want %d", rec.Code, tt.wantCode)
				}
			})
		}
	})

	t.Run("absolute time parsing", func(t *testing.T) {
		tests := []struct {
			name        string
			content     string
			layout      string
			wantYear    int
			wantMonth   time.Month
			wantDay     int
			wantHour    int
			wantMinute  int
			wantSecond  int
			shouldParse bool
		}{
			{
				name:        "RFC3339 format",
				content:     "2023-06-15T10:30:00Z",
				layout:      time.RFC3339,
				wantYear:    2023,
				wantMonth:   time.June,
				wantDay:     15,
				wantHour:    10,
				wantMinute:  30,
				wantSecond:  0,
				shouldParse: true,
			},
			{
				name:        "RFC3339 with @ prefix",
				content:     "@2023-06-15T10:30:00Z",
				layout:      time.RFC3339,
				wantYear:    2023,
				wantMonth:   time.June,
				wantDay:     15,
				wantHour:    10,
				wantMinute:  30,
				wantSecond:  0,
				shouldParse: true,
			},
			{
				name:        "DateOnly format",
				content:     "2023-06-15",
				layout:      time.DateOnly,
				wantYear:    2023,
				wantMonth:   time.June,
				wantDay:     15,
				wantHour:    0,
				wantMinute:  0,
				wantSecond:  0,
				shouldParse: true,
			},
			{
				name:        "DateTime format with space",
				content:     "2023-06-15 10:30:00",
				layout:      time.DateTime,
				wantYear:    2023,
				wantMonth:   time.June,
				wantDay:     15,
				wantHour:    10,
				wantMinute:  30,
				wantSecond:  0,
				shouldParse: true,
			},
			{
				name:        "invalid time format",
				content:     "not-a-valid-time",
				layout:      time.RFC3339,
				shouldParse: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}

				var capturedTime time.Time
				var handlerCalled bool
				handler := Middleware(filePath, tt.layout)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						handlerCalled = true
						capturedTime = time.Now()
						w.WriteHeader(http.StatusOK)
					}),
				)

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if !handlerCalled {
					t.Fatal("handler was not called")
				}

				if !tt.shouldParse {
					return // For invalid format, just verify handler was called
				}

				// Compare in UTC to avoid timezone issues
				utc := capturedTime.UTC()
				if utc.Year() != tt.wantYear {
					t.Errorf("year: got %d, want %d", utc.Year(), tt.wantYear)
				}
				if utc.Month() != tt.wantMonth {
					t.Errorf("month: got %v, want %v", utc.Month(), tt.wantMonth)
				}
				if utc.Day() != tt.wantDay {
					t.Errorf("day: got %d, want %d", utc.Day(), tt.wantDay)
				}
				if utc.Hour() != tt.wantHour {
					t.Errorf("hour: got %d, want %d", utc.Hour(), tt.wantHour)
				}
				if utc.Minute() != tt.wantMinute {
					t.Errorf("minute: got %d, want %d", utc.Minute(), tt.wantMinute)
				}
				if utc.Second() != tt.wantSecond {
					t.Errorf("second: got %d, want %d", utc.Second(), tt.wantSecond)
				}
			})
		}
	})

	t.Run("offset duration", func(t *testing.T) {
		tests := []struct {
			name           string
			content        string
			expectedOffset time.Duration
			shouldParse    bool
		}{
			{
				name:           "positive offset +1h",
				content:        "+1h",
				expectedOffset: time.Hour,
				shouldParse:    true,
			},
			{
				name:           "positive offset +30m",
				content:        "+30m",
				expectedOffset: 30 * time.Minute,
				shouldParse:    true,
			},
			{
				name:           "negative offset -30m",
				content:        "-30m",
				expectedOffset: -30 * time.Minute,
				shouldParse:    true,
			},
			{
				name:           "negative offset -2h",
				content:        "-2h",
				expectedOffset: -2 * time.Hour,
				shouldParse:    true,
			},
			{
				name:           "complex offset +1h30m",
				content:        "+1h30m",
				expectedOffset: time.Hour + 30*time.Minute,
				shouldParse:    true,
			},
			{
				name:        "invalid duration",
				content:     "+invalid",
				shouldParse: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}

				realNow := time.Now()
				var capturedTime time.Time
				var handlerCalled bool
				handler := Middleware(filePath, time.RFC3339)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						handlerCalled = true
						capturedTime = time.Now()
						w.WriteHeader(http.StatusOK)
					}),
				)

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if !handlerCalled {
					t.Fatal("handler was not called")
				}

				if !tt.shouldParse {
					return
				}

				diff := capturedTime.Sub(realNow)
				tolerance := time.Second
				if diff < tt.expectedOffset-tolerance || diff > tt.expectedOffset+tolerance {
					t.Errorf("offset: got %v, want %v (±%v)", diff, tt.expectedOffset, tolerance)
				}
			})
		}
	})

	t.Run("time progression options", func(t *testing.T) {
		tests := []struct {
			name         string
			content      string
			layout       string
			wantBaseTime time.Time
			shouldParse  bool
		}{
			{
				name:         "with + option (ratio 1.0)",
				content:      "2023-06-15 10:30:00 +",
				layout:       time.DateTime,
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:         "with x1 option",
				content:      "2023-06-15 10:30:00 x1",
				layout:       time.DateTime,
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:         "with x2 option",
				content:      "2023-06-15 10:30:00 x2",
				layout:       time.DateTime,
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:         "with x0.5 option",
				content:      "2023-06-15 10:30:00 x0.5",
				layout:       time.DateTime,
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:        "x without ratio",
				content:     "2023-06-15 10:30:00 x",
				layout:      time.DateTime,
				shouldParse: false,
			},
			{
				name:        "invalid ratio xabc",
				content:     "2023-06-15 10:30:00 xabc",
				layout:      time.DateTime,
				shouldParse: false,
			},
			{
				name:        "zero ratio x0",
				content:     "2023-06-15 10:30:00 x0",
				layout:      time.DateTime,
				shouldParse: false,
			},
			{
				name:        "negative ratio x-1",
				content:     "2023-06-15 10:30:00 x-1",
				layout:      time.DateTime,
				shouldParse: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}

				var capturedTime time.Time
				var handlerCalled bool
				handler := Middleware(filePath, tt.layout)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						handlerCalled = true
						capturedTime = time.Now()
						w.WriteHeader(http.StatusOK)
					}),
				)

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if !handlerCalled {
					t.Fatal("handler was not called")
				}

				if !tt.shouldParse {
					return
				}

				// Check that captured time is close to the expected base time
				// (allowing small tolerance for execution time)
				diff := capturedTime.UTC().Sub(tt.wantBaseTime)
				if diff < 0 {
					diff = -diff
				}
				tolerance := time.Second
				if diff > tolerance {
					t.Errorf("time: got %v, want %v (±%v)", capturedTime.UTC(), tt.wantBaseTime, tolerance)
				}
			})
		}
	})

	t.Run("RFC3339 with options", func(t *testing.T) {
		tests := []struct {
			name         string
			content      string
			wantBaseTime time.Time
			shouldParse  bool
		}{
			{
				name:         "RFC3339 with + option",
				content:      "2023-06-15T10:30:00Z +",
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:         "RFC3339 with x2 option",
				content:      "2023-06-15T10:30:00Z x2",
				wantBaseTime: time.Date(2023, 6, 15, 10, 30, 0, 0, time.UTC),
				shouldParse:  true,
			},
			{
				name:    "RFC3339 with invalid option",
				content: "2023-06-15T10:30:00Z invalid",
				// This should not be recognized as having an option
				// because "invalid" doesn't start with "x" or equal "+"
				// so body becomes the entire string and parsing fails
				shouldParse: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "time.txt")
				if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}

				var capturedTime time.Time
				var handlerCalled bool
				handler := Middleware(filePath, time.RFC3339)(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						handlerCalled = true
						capturedTime = time.Now()
						w.WriteHeader(http.StatusOK)
					}),
				)

				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)

				if !handlerCalled {
					t.Fatal("handler was not called")
				}

				if !tt.shouldParse {
					return
				}

				diff := capturedTime.UTC().Sub(tt.wantBaseTime)
				if diff < 0 {
					diff = -diff
				}
				tolerance := time.Second
				if diff > tolerance {
					t.Errorf("time: got %v, want %v (±%v)", capturedTime.UTC(), tt.wantBaseTime, tolerance)
				}
			})
		}
	})
}
