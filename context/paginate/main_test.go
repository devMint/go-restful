package paginate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devMint/go-restful/request"
	"github.com/devMint/go-restful/response"
	"github.com/stretchr/testify/assert"
)

func Test_PaginateContext_WithParams(t *testing.T) {
	handler := request.HandleAction(paginationHandler)
	handlerToTest := Paginate(30, 0)(handler)

	request, _ := http.NewRequest("GET", "/?take=12&skip=3", nil)
	request.Header.Set("content-type", "application/json")

	response := httptest.NewRecorder()
	handlerToTest.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assert.Equal(t, "{\"data\":{\"skip\":3,\"take\":12}}", response.Body.String())
}

func Test_PaginateContext_EmptyParams(t *testing.T) {
	handler := request.HandleAction(paginationHandler)
	handlerToTest := Paginate(30, 0)(handler)

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")

	response := httptest.NewRecorder()
	handlerToTest.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assert.Equal(t, "{\"data\":{\"skip\":0,\"take\":30}}", response.Body.String())
}

func Test_PaginateContext_InvalidParam(t *testing.T) {
	handler := request.HandleAction(paginationHandler)
	handlerToTest := Paginate(30, 0)(handler)

	request, _ := http.NewRequest("GET", "/?take=a", nil)
	request.Header.Set("content-type", "application/json")

	response := httptest.NewRecorder()
	handlerToTest.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
}

func paginationHandler(req request.Request) response.Response {
	ctx := req.Context()
	take, skip := ctx.Value(PaginateTake).(int), ctx.Value(PaginateSkip).(int)

	return response.Ok(map[string]int{"take": take, "skip": skip})
}
