package errs

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/response"
	"github.com/go-chi/render"
)

var (
	ErrDuplicateAlias = errors.New("duplicate alias")
	ErrFailendToSend  = errors.New("failed to send")
	ErrInvalidStruct  = errors.New("invalid struct")
	ErrInvalidUrl     = errors.New("invalid url")
)

func ErrWithReq(log *slog.Logger, r *http.Request, w http.ResponseWriter, err error, msgLog string, msgRes string, status int) {
	log.Error(msgLog, slog.Any("error", err))
	render.Status(r, status)
	render.JSON(w, r, response.Error(msgRes))
}
