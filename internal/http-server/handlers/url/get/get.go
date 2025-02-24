package get

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/errs"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/response"
	"github.com/go-chi/chi/v5/middleware"
)

type Request struct {
	Alias string `json:"alias" valid:"required"`
}

type Response struct {
	response.Response
	Url string `json:"URL"`
}

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(logger *slog.Logger, getter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.URL.get.New"
		log := logger.With(slog.String("operation", operation), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errs.ErrWithReq(log, r, w, err, "failed to decode request", "internal error", http.StatusBadRequest)
			return
		}
		// check validate
		log.Info("request decoded", slog.Any("request", req))
		if err := ValidateReq(req); err != nil {
			errs.ErrWithReq(log, r, w, err, "failed to validate request", "porvide alias", http.StatusBadRequest)
			return
		}
		log.Info("request validated", slog.Any("request", req))
		log.Info("getting URL", slog.String("alias", req.Alias))

		url, err := getter.GetUrl(req.Alias)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errs.ErrWithReq(log, r, w, err, "URL not found", "not found", http.StatusNotFound)
				return
			}
			errs.ErrWithReq(log, r, w, err, "failed to get URL", "internal error", http.StatusInternalServerError)
			return
		}
		log.Info("URL got", slog.String("URL", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func ValidateReq(req Request) error {
	if req.Alias == "" {
		return errs.ErrInvalidStruct
	}
	return nil
}
