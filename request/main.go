package request

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"gitlab.com/devmint/go-restful/response"
)

var validate RequestBodyValidation

// RestfulHandler replacement for http.HandlerFunc
type RestfulHandler func(Request) response.Response

// ContextHandler special replacement for http.Handler. It's should not be use for middlewares
// but more with common actions between restful handlers.
type ContextHandler func(Request) (context.Context, response.Response)

// Request replacement for *http.Request with easy access to most important variables or parameters.
type Request interface {
	Param(key string) string
	Query(key string, onMissing ...string) string
	Body(typeOfBody interface{}) error
	Context() context.Context
	Request() *http.Request
}

type RequestBodyValidation interface {
	Struct(s interface{}) error
}

const (
	appJSON = "application/json"
	appXML  = "application/xml"
)

// HandleAction replacement for http.HandlerFunc
func HandleAction(cb func(req Request) response.Response) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := cb(wrapRequest(r))
		for name, value := range response.Header() {
			w.Header().Set(name, value)
		}

		switch r.Header.Get("content-type") {
		case appJSON:
			renderJSON(w, response)
		case appXML:
			renderXML(w, response)
		default:
			renderJSON(w, response)
		}
	})
}

// HandleContext replacement for func(http.Handler) http.Handler
func HandleContext(cb func(req Request) (context.Context, response.Response)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, res := cb(wrapRequest(r))
			if res != nil {
				switch r.Header.Get("content-type") {
				case appJSON:
					renderJSON(w, res)
				case appXML:
					renderXML(w, res)
				default:
					renderJSON(w, res)
				}
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type nativeRequest struct {
	request   *http.Request
	validator RequestBodyValidation
}

func (r nativeRequest) Request() *http.Request { return r.request }

func (r nativeRequest) Context() context.Context { return r.request.Context() }

func (r nativeRequest) Param(key string) string { return chi.URLParam(r.request, key) }

func (r nativeRequest) Query(key string, onMissing ...string) string {
	keys, ok := r.request.URL.Query()[key]
	if !ok || len(keys) == 0 {
		if len(onMissing) > 0 {
			return onMissing[0]
		}
		return ""
	}

	return keys[0]
}

func (r nativeRequest) Body(typeOfBody interface{}) error {
	body := r.request.Body
	if body == nil {
		return errors.New("empty body from request")
	}

	contentType := appJSON
	if r.request.Header.Get("content-type") != "" {
		contentType = r.request.Header.Get("content-type")
	}

	var err error
	switch contentType {
	case appJSON:
		err = json.NewDecoder(body).Decode(&typeOfBody)
	case appXML:
		err = xml.NewDecoder(body).Decode(&typeOfBody)
	default:
		err = fmt.Errorf("content type '%s' is unsupoorted", r.request.Header.Get("content-type"))
	}

	if err == nil {
		err = r.validator.Struct(typeOfBody)
	}

	return err
}

func wrapRequest(r *http.Request) nativeRequest {
	return nativeRequest{
		request:   r,
		validator: validate,
	}
}

func init() {
	validate = validator.New()
}

func RegisterValidator(v RequestBodyValidation) {
	validate = v
}

func renderJSON(w http.ResponseWriter, r response.Response) {
	w.WriteHeader(r.StatusCode())
	w.Header().Set("content-type", appJSON)
	w.Write([]byte(r.GetJSON()))
}

func renderXML(w http.ResponseWriter, r response.Response) {
	w.WriteHeader(r.StatusCode())
	w.Header().Set("content-type", appXML)
	w.Write([]byte(r.GetXML()))
}
