package faketimehttp

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/akm/time/faketime"
)

func Middleware(filePath string, layout string) func(next http.Handler) http.Handler {
	provider := faketime.NewFileProvider(filePath)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			s, err := provider.Get(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get faketime", "error", err, "file", filePath)
				http.Error(w, "faketime error", http.StatusInternalServerError)
				return
			}

			if s == "" {
				next.ServeHTTP(w, r)
				return
			}

			ft, err := faketime.Parse(s, layout)
			if err != nil {
				slog.ErrorContext(ctx, "failed to parse faketime", "error", err, "file", filePath, "content", s)
				http.Error(w, "faketime error", http.StatusInternalServerError)
				return
			}

			_ = ft.Run(ctx, func(ctx context.Context) error {
				next.ServeHTTP(w, r)
				return nil
			})
		})
	}
}
