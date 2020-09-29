package restful

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devMint/go-restful/request"
	"github.com/devMint/go-restful/response"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func Test_GetRoute(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Get("/", func(r request.Request) response.Response {
		return response.Ok("test")
	})

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":\"test\"}", response.Body.String())
}

func Test_GetRoute_WithContext(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.With(func(r request.Request) (context.Context, response.Response) {
		return context.WithValue(r.Context(), "test", "test2"), nil
	}).Get("/", func(r request.Request) response.Response {
		return response.Ok(r.Context().Value("test"))
	})

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":\"test2\"}", response.Body.String())
}

func Test_GroupRoute(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Group(func(r Router) {
		router.Get("/", func(request.Request) response.Response { return response.Ok("ok") })
	})

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":\"ok\"}", response.Body.String())
}

func Test_RouteNested(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Route("/path", func(r Router) {
		r.Get("/a", func(request.Request) response.Response { return response.Ok("a") })
		r.Post("/b", func(request.Request) response.Response { return response.Ok("b") })
	})

	// first request
	responseA := httptest.NewRecorder()
	requestA, _ := http.NewRequest("GET", "/path/a", nil)
	router.ServeHTTP(responseA, requestA)

	assert.Equal(t, "{\"data\":\"a\"}", responseA.Body.String())

	// second request
	responseB := httptest.NewRecorder()
	requestB, _ := http.NewRequest("POST", "/path/b", nil)
	router.ServeHTTP(responseB, requestB)

	assert.Equal(t, "{\"data\":\"b\"}", responseB.Body.String())
}

func Benchmark_GetRoute(b *testing.B) {
	router := NewRouter(chi.NewMux())
	router.Get("/", func(r request.Request) response.Response {
		return response.Ok("test")
	})

	for i := 0; i < b.N; i++ {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(response, request)
	}
}

func Benchmark_GetRoute_WithContext(b *testing.B) {
	router := NewRouter(chi.NewMux())
	router.With(func(r request.Request) (context.Context, response.Response) {
		return context.WithValue(r.Context(), "test", "test2"), nil
	}).Get("/", func(r request.Request) response.Response {
		return response.Ok(r.Context().Value("test"))
	})

	for i := 0; i < b.N; i++ {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(response, request)
	}
}
