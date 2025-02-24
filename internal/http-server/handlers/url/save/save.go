package save

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/errs"
	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/response"
	"github.com/DblMOKRQ/Url-Shortener/internal/lib/random"

	// "github.com/asaskevich/govalidator"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// TODO:
// 4. Может быть такое что при генерации алиаса уже есть алиас с таким именем
const (
	lengthAlias    = 6
	ErrAliasExists = "duplicate value in alias"
)

type Request struct {
	Url   string `json:"URL" valid:"required"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) error
}

func New(logger *slog.Logger, saver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "handlers.URL.save.New"
		log := logger.With(slog.String("operation", operation), slog.String("request_id", middleware.GetReqID(r.Context())))
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errs.ErrWithReq(log, r, w, err, "failed to decode request", "internal error", http.StatusBadRequest)
			return
		}
		log.Info("request decoded", slog.Any("request", req))
		//check validate
		if err := ValidateReq(req); err != nil {
			switch err {
			case errs.ErrInvalidStruct:
				errs.ErrWithReq(log, r, w, err, "the URL is not specified", "provide URL", http.StatusBadRequest)
			case errs.ErrInvalidUrl:
				errs.ErrWithReq(log, r, w, err, "URL is not valid", "URL not valid", http.StatusBadRequest)
			default:
				errs.ErrWithReq(log, r, w, err, "invalid request", "internal error", http.StatusBadRequest)
			}
			return
		}
		//
		// genetate alias
		alias := generatAlias(req.Alias)
		// save URL
		err := saver.SaveUrl(req.Url, alias)
		if err != nil {
			if err == errs.ErrDuplicateAlias {
				errs.ErrWithReq(log, r, w, err, "alias already exists", "use another alias", http.StatusInternalServerError)
				return
			}
			errs.ErrWithReq(log, r, w, err, "failed to save URL", "internal error", http.StatusInternalServerError)
			return
		}
		log.Info("URL saved", slog.String("alias", alias), slog.String("URL", req.Url))
		//send response
		render.JSON(w, r, Response{Response: response.Ok(), Alias: alias})

	}
}
func generatAlias(alias string) string {
	if alias == "" {
		return random.RandomString(lengthAlias)
	}
	return alias
}
func ValidateReq(req Request) error {
	if req.Url == "" {
		return errs.ErrInvalidStruct
	}
	if u, err := url.Parse(req.Url); !(err == nil && u.Scheme != "" && u.Host != "") {
		return errs.ErrInvalidUrl
	}
	return nil
}
