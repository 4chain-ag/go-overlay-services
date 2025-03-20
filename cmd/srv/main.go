package main

import (
	"net/http"

	"github.com/4chain-ag/go-overlay-services/pkg/server"
	"github.com/gookit/slog"
)

func main() {
	opts := []server.HTTPOption{
		server.WithConfig(&server.Config{
			AdminBearerToken: "12345678secret!",
			Addr:             "localhost",
			Port:             8080,
		}),
		server.WithLoggingMiddleware(loggingMiddleware),
	}

	httpAPI := server.New(opts...)
	httpAPI.ListenAndServe()
}

func loggingMiddleware(next http.Handler) http.Handler {
	slog.SetLogLevel(slog.DebugLevel)
	slog.SetFormatter(slog.NewJSONFormatter(func(f *slog.JSONFormatter) {
		f.PrettyPrint = true
	}))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.WithFields(slog.M{
			"category":    "service",
			"method":      r.Method,
			"remote-addr": r.RemoteAddr,
			"request-uri": r.RequestURI,
		})
		logger.Info("log-line")
		next.ServeHTTP(w, r)
	})
}
