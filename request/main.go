package request

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/devMint/go-restful/response"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

type RestfulHandler func(Request) response.Response
type ContextHandler func(req Request) (context.Context, response.Response)

type Request interface {
	Param(key string) string
	Query(key string, onMissing ...string) string
	Body(typeOfBody interface{}) error
	BodyWithValidation(typeOfBody interface{}) error
	Context() context.Context
}

const (
	appJSON = "application/json"
	appXML  = "application/xml"
)

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
	validator *validator.Validate
}

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

	var err error
	switch r.request.Header.Get("content-type") {
	case appJSON:
		json.NewDecoder(body).Decode(&typeOfBody)
	case appXML:
		xml.NewDecoder(body).Decode(&typeOfBody)
	default:
		err = fmt.Errorf("content type '%s' is unsupoorted", r.request.Header.Get("content-type"))
	}

	return err
}

func (r nativeRequest) BodyWithValidation(typeOfBody interface{}) error {
	if err := r.Body(&typeOfBody); err != nil {
		return err
	}

	return r.validator.Struct(typeOfBody)
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
