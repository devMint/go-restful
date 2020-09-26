package request

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devMint/go-restful/response"
	"github.com/stretchr/testify/assert"
)

const (
	responseInJSON = "{\"type\":\"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html\",\"title\":\"Not Found\",\"detail\":\"employee not found\",\"status\":404}"
	responseInXML  = "<response><type>http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html</type><title>Not Found</title><detail>employee not found</detail><status>404</status></response>"
)

var errFoo = errors.New("employee not found")

func Test_ErrorResponse_JSON(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(notFoundHandler))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assert.Equal(t, responseInJSON, response.Body.String())
}

func Test_ErrorResponse_XML(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(notFoundHandler))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/xml")
	response := httpResponse(handler, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "application/xml", response.Header().Get("content-type"))
	assert.Equal(t, responseInXML, response.Body.String())
}

func Test_OkResponse_Json(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(collectionHandler))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assert.Equal(t, "{\"data\":[\"a\",\"b\",\"c\"]}", response.Body.String())
}

func Test_Context_Valid(t *testing.T) {
	handler := HandleAction(collectionHandler)
	handlerToTest := HandleContext(validContext)(handler)

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")

	response := httptest.NewRecorder()
	handlerToTest.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
	assert.Equal(t, "{\"data\":[\"d\",\"e\",\"f\"]}", response.Body.String())
}

func Test_Context_InValid(t *testing.T) {
	handler := HandleAction(collectionHandler)
	handlerToTest := HandleContext(invalidContext)(handler)

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")

	response := httptest.NewRecorder()
	handlerToTest.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("content-type"))
}

func collectionHandler(r Request) response.Response {
	ctx := r.Context()
	if ctx.Value("valid") == nil {
		return response.Ok("a", "b", "c")
	}

	return response.Ok(ctx.Value("valid"))
}

func notFoundHandler(r Request) response.Response {
	return response.NotFound(errFoo)
}

func validContext(r Request) (context.Context, response.Response) {
	return context.WithValue(r.Context(), "valid", []string{"d", "e", "f"}), nil
}

func invalidContext(r Request) (context.Context, response.Response) {
	return nil, response.NotFound(errors.New("context not found"))
}

func httpResponse(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	handler(w, req)
	return w
}
