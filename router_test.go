package restful

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	router.Get("/", func(r request.Request) response.Response { return response.Ok("test") })

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":\"test\"}", response.Body.String())
}

func Test_GetRoute_Param(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Get("/{id}", func(r request.Request) response.Response { return response.Ok(r.Param("id")) })

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/12", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":\"12\"}", response.Body.String())
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

func Test_PostRoute_DefaultValidator(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Post("/", func(r request.Request) response.Response {
		body := customBody{}
		if err := r.Body(&body); err != nil {
			return response.BadRequest(err)
		}

		return response.Ok(body.A)
	})

	body, _ := json.Marshal(map[string]int{"a": 16})
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "{\"data\":16}", response.Body.String())
}

func Test_PostRoute_CustomValidator(t *testing.T) {
	router := NewRouter(chi.NewMux(), RouterOptions{Validator: customValidator{}})
	router.Post("/", func(r request.Request) response.Response {
		body := customBody{}
		if err := r.Body(&body); err != nil {
			return response.BadRequest(err)
		}

		return response.Ok(body.A)
	})

	body, _ := json.Marshal(map[string]int{"a": 16})
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "some-random-error")
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

func Test_RouteNested_SameSlash(t *testing.T) {
	router := NewRouter(chi.NewMux())
	router.Get("/", func(request.Request) response.Response { return response.Ok("a") })
	router.Route("/path", func(r Router) {
		r.Get("/", func(request.Request) response.Response { return response.Ok("b") })
	})

	// first request
	responseA := httptest.NewRecorder()
	requestA, _ := http.NewRequest("GET", "/path", nil)
	router.ServeHTTP(responseA, requestA)

	assert.Equal(t, "{\"data\":\"b\"}", responseA.Body.String())

	// second request
	responseB := httptest.NewRecorder()
	requestB, _ := http.NewRequest("POST", "/", nil)
	router.ServeHTTP(responseB, requestB)

	assert.Equal(t, "{\"data\":\"a\"}", responseB.Body.String())
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

func Benchmark_PostRoute_Validation(b *testing.B) {
	router := NewRouter(chi.NewMux())
	router.Post("/", func(r request.Request) response.Response {
		body := customBody{}
		if err := r.Body(&body); err != nil {
			return response.BadRequest(err)
		}

		return response.Ok(body.A)
	})

	for i := 0; i < b.N; i++ {
		body, _ := json.Marshal(map[string]int{"a": i})

		response := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/", bytes.NewReader(body))
		router.ServeHTTP(response, request)
	}
}

type customValidator struct{}

func (v customValidator) Struct(s interface{}) error { return errors.New("some-random-error") }

type customBody struct {
	A int `json:"a" xml:"a" validate:"required,gte=13"`
}
