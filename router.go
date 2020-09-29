package restful

import (
	"net/http"

	"github.com/devMint/go-restful/request"
	"github.com/go-chi/chi"
)

type Router interface {
	http.Handler

	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...request.ContextHandler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...request.ContextHandler) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a sub-Router along a `pattern`` string.
	Route(pattern string, fn func(r Router)) Router

	Mount(pattern string, h http.Handler)

	// HTTP-method routing along `pattern`
	Connect(pattern string, h request.RestfulHandler)
	Delete(pattern string, h request.RestfulHandler)
	Get(pattern string, h request.RestfulHandler)
	Head(pattern string, h request.RestfulHandler)
	Options(pattern string, h request.RestfulHandler)
	Patch(pattern string, h request.RestfulHandler)
	Post(pattern string, h request.RestfulHandler)
	Put(pattern string, h request.RestfulHandler)
	Trace(pattern string, h request.RestfulHandler)
}

type restfulRouter struct {
	r chi.Router
}

func NewRouter(plainRouter chi.Router) Router {
	return restfulRouter{r: plainRouter}
}

func (router restfulRouter) Use(middlewares ...request.ContextHandler) {
	var httpMiddlewares []func(http.Handler) http.Handler
	for _, middleware := range middlewares {
		httpMiddlewares = append(httpMiddlewares, request.HandleContext(middleware))
	}

	router.r.Use(httpMiddlewares...)
}

func (router restfulRouter) With(middlewares ...request.ContextHandler) Router {
	var httpMiddlewares []func(http.Handler) http.Handler
	for _, middleware := range middlewares {
		httpMiddlewares = append(httpMiddlewares, request.HandleContext(middleware))
	}

	return restfulRouter{r: router.r.With(httpMiddlewares...)}
}

func (router restfulRouter) Group(fn func(r Router)) Router {
	newRouter := router.With()
	if fn != nil {
		fn(newRouter)
	}
	return newRouter
}

func (router restfulRouter) Route(pattern string, fn func(r Router)) Router {
	newRouter := restfulRouter{r: chi.NewMux()}
	if fn != nil {
		fn(newRouter)
	}

	router.Mount(pattern, newRouter)
	return newRouter
}

func (router restfulRouter) Mount(pattern string, h http.Handler) {
	router.r.Mount(pattern, h)
}

func (router restfulRouter) method(method string, pattern string, h request.RestfulHandler) {
	router.r.MethodFunc(method, pattern, request.HandleAction(h))
}

func (router restfulRouter) Connect(pattern string, h request.RestfulHandler) {
	router.method(http.MethodConnect, pattern, h)
}

func (router restfulRouter) Delete(pattern string, h request.RestfulHandler) {
	router.method(http.MethodDelete, pattern, h)
}

func (router restfulRouter) Get(pattern string, h request.RestfulHandler) {
	router.method(http.MethodGet, pattern, h)
}

func (router restfulRouter) Head(pattern string, h request.RestfulHandler) {
	router.method(http.MethodHead, pattern, h)
}

func (router restfulRouter) Options(pattern string, h request.RestfulHandler) {
	router.method(http.MethodOptions, pattern, h)
}

func (router restfulRouter) Patch(pattern string, h request.RestfulHandler) {
	router.method(http.MethodPatch, pattern, h)
}

func (router restfulRouter) Post(pattern string, h request.RestfulHandler) {
	router.method(http.MethodPost, pattern, h)
}

func (router restfulRouter) Put(pattern string, h request.RestfulHandler) {
	router.method(http.MethodPut, pattern, h)
}

func (router restfulRouter) Trace(pattern string, h request.RestfulHandler) {
	router.method(http.MethodTrace, pattern, h)
}

func (router restfulRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.r.ServeHTTP(w, r)
}
