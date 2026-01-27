package faketimehttp

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	orig "time"

	"github.com/akm/time"
	"github.com/akm/time/testtime"
)

func Middleware(filePath string, layout string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			stat, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					slog.DebugContext(ctx, "time file does not exist, proceeding without setting time", "file", filePath)
				} else {
					slog.WarnContext(ctx, "failed to stat file", "error", err, "file", filePath)
				}
				next.ServeHTTP(w, r)
				return
			}
			if stat.IsDir() {
				slog.ErrorContext(ctx, "path is a directory, not a file", "file", filePath)
				next.ServeHTTP(w, r)
				return
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				slog.ErrorContext(ctx, "failed to read file", "error", err, "file", filePath)
				next.ServeHTTP(w, r)
				return
			}
			s := strings.TrimSpace(string(data))
			if s == "" {
				slog.DebugContext(ctx, "time file is empty, proceeding without setting time", "file", filePath)
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(s, " ")
			var body string
			var opt string
			if len(parts) > 1 {
				last := parts[len(parts)-1]
				if last == "+" || strings.HasPrefix(last, "x") {
					opt = last
					body = strings.Join(parts[:len(parts)-1], " ")
				} else {
					body = s
				}
			} else {
				body = s
			}

			var t time.Time
			if strings.HasPrefix(body, "+") || strings.HasPrefix(body, "-") {
				offsetDuration, err := time.ParseDuration(body)
				if err != nil {
					slog.ErrorContext(ctx, "failed to parse duration from file content", "error", err, "file", filePath, "content", s)
					next.ServeHTTP(w, r)
					return
				}
				t = time.Now().Add(offsetDuration)
			} else {
				body = strings.TrimPrefix(body, "@")
				var err error
				t, err = time.Parse(layout, body)
				if err != nil {
					slog.ErrorContext(ctx, "failed to parse time from file content", "error", err, "file", filePath, "content", s)
					next.ServeHTTP(w, r)
					return
				}
			}

			if opt == "" {
				defer testtime.SetTime(&t)()
				next.ServeHTTP(w, r)
				return
			}

			var ratioStr string
			if strings.HasPrefix(opt, "x") {
				ratioStr = strings.TrimPrefix(opt, "x")
				if ratioStr == "" {
					slog.ErrorContext(ctx, "missing ratio after 'x' in file content", "file", filePath, "content", s)
					next.ServeHTTP(w, r)
					return
				}
			}

			var ratio float64
			if ratioStr == "" {
				ratio = 1.0
			} else {
				var err error
				ratio, err = strconv.ParseFloat(ratioStr, 64)
				if err != nil {
					slog.ErrorContext(ctx, "failed to parse ratio from file content", "error", err, "file", filePath, "content", s)
					next.ServeHTTP(w, r)
					return
				}
				if ratio <= 0 {
					slog.ErrorContext(ctx, "ratio must be positive", "file", filePath, "content", s)
					next.ServeHTTP(w, r)
					return
				}
			}

			t0 := orig.Now()
			defer testtime.SetTimeFunc(func() time.Time {
				elapsed := time.Duration(float64(orig.Since(t0)) * ratio)
				return t.Add(elapsed)
			})()
			next.ServeHTTP(w, r)
		})
	}
}
