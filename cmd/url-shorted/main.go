package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/DblMOKRQ/Url-Shortener/internal/config"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/handlers/url/delete"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/handlers/url/get"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/handlers/url/save"
	postgresql "github.com/DblMOKRQ/Url-Shortener/internal/storage/postgreSQL"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	storage, err := postgresql.New(cfg.User, cfg.Password, cfg.DBName, cfg.Sslmode)
	if err != nil {
		log.Error("failed to init storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer storage.Close()
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer) // Восстановление после паники

	router.Route("/", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			"myuser": "mypassword",
		}))
		r.Post("/saveurl", save.New(log, storage))
		r.Post("/deleteurl", delete.New(log, storage))
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	router.Post("/geturl", get.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Addres))

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.Any("error", err))
	}
	log.Info("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case localEnv:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case devEnv:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prodEnv:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return logger

}
