package delete

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/errs"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/response"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias" valid:"required"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type UrlDeleted interface {
	DeleteUrl(alias string) error
}

func New(logger *slog.Logger, deleted UrlDeleted) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.URL.delete.New"
		log := logger.With(slog.String("operation", operation), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errs.ErrWithReq(log, r, w, err, "failed to decode request", "internal error", http.StatusBadRequest)
			return
		}
		// check validate
		log.Info("request decoded", slog.Any("request", req))
		if err := ValidateReq(req); err != nil {
			errs.ErrWithReq(log, r, w, err, "failed to validate request", "provide alias", http.StatusBadRequest)
			return
		}
		log.Info("request validated", slog.Any("request", req))

		err := deleted.DeleteUrl(req.Alias)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errs.ErrWithReq(log, r, w, err, "URL not found", "not found", http.StatusNotFound)
				return
			}
			errs.ErrWithReq(log, r, w, err, "failed to delete URL", "internal error", http.StatusInternalServerError)
			return
		}
		log.Info("URL deleted", slog.String("alias", req.Alias))
		render.JSON(w, r, Response{Response: response.Ok(), Alias: req.Alias})

	}
}
func ValidateReq(req Request) error {
	if req.Alias == "" {
		return errs.ErrInvalidStruct
	}
	return nil
}
