package save_test

import (
	"testing"

	"github.com/DblMOKRQ/Url-Shortener/internal/http-server/api/errs"
	. "github.com/DblMOKRQ/Url-Shortener/internal/http-server/handlers/url/save"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name   string
		req    Request
		expErr error
	}{
		{
			name: "empty url",
			req: Request{
				Url: "",
			},
			expErr: errs.ErrInvalidStruct,
		},
		{
			name: "invalid url",
			req: Request{
				Url: "**********",
			},
			expErr: errs.ErrInvalidUrl,
		},
		{
			name: "valid url",
			req: Request{
				Url: "http://valid",
			},
			expErr: nil,
		},
		{
			name: "valid url with alias",
			req: Request{
				Url:   "http://valid",
				Alias: "alias",
			},
			expErr: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateReq(c.req)
			if err != c.expErr {
				t.Errorf("expected error %v, got %v", c.expErr, err)
			}
		})
	}
}
