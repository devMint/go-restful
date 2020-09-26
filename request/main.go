package request

import (
	"context"
	"net/http"

	"github.com/devMint/go-restful/response"
	"github.com/go-chi/chi"
)

type RestfulHandler func(Request) response.Response
type ContextHandler func(req Request) (context.Context, response.Response)

type Request interface {
	Param(key string) string
	Query(key string, onMissing string) string
	Context() context.Context
}

const (
	appJSON = "application/json"
	appXML  = "application/xml"
)

func HandleAction(cb func(req Request) response.Response) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := cb(nativeRequest{request: r})

		switch r.Header.Get("content-type") {
		case appJSON:
			renderJson(w, response)
		case appXML:
			renderXml(w, response)
		default:
			renderJson(w, response)
		}
	})
}

func HandleContext(cb func(req Request) (context.Context, response.Response)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, res := cb(nativeRequest{request: r})
			if res != nil {
				switch r.Header.Get("content-type") {
				case appJSON:
					renderJson(w, res)
				case appXML:
					renderXml(w, res)
				default:
					renderJson(w, res)
				}
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type nativeRequest struct {
	request *http.Request
}

func (r nativeRequest) Context() context.Context { return r.request.Context() }
func (r nativeRequest) Param(key string) string  { return chi.URLParam(r.request, key) }
func (r nativeRequest) Query(key string, onMissing string) string {
	keys, ok := r.request.URL.Query()[key]
	if !ok || len(keys) == 0 {
		return onMissing
	}

	return keys[0]
}

func renderJson(w http.ResponseWriter, r response.Response) {
	w.WriteHeader(r.StatusCode())
	w.Header().Set("content-type", appJSON)
	w.Write([]byte(r.GetJSON()))
}

func renderXml(w http.ResponseWriter, r response.Response) {
	w.WriteHeader(r.StatusCode())
	w.Header().Set("content-type", appXML)
	w.Write([]byte(r.GetXML()))
}
