package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func Test_CustomHeaders(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(customHeaders))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	assert.Equal(t, "dolor-sit-amet", response.Header().Get("lorem-ipsum"))
}

func Test_ParseBody(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(bodyToResponse))
	body, _ := json.Marshal(map[string]string{"a": "lorem-ipsum"})

	request, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	assert.Equal(t, "{\"data\":{\"a\":\"lorem-ipsum\"}}", response.Body.String())
}

func Test_ParseBody_WithValidation(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(bodyToResponseWithValidation))
	body, _ := json.Marshal(map[string]string{"a": "lorem-ipsum"})

	request, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	assert.Contains(t, response.Body.String(), "Key: 'customBody.A' Error:Field validation for 'A' failed on the 'iscolor' tag")
}

func Test_Redirect(t *testing.T) {
	handler := http.HandlerFunc(HandleAction(redirectHandler))

	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Set("content-type", "application/json")
	response := httpResponse(handler, request)

	fmt.Printf("%+v", response.Header())

	assert.Equal(t, "http://www.onet.pl", response.Header().Get("Location"))
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

func customHeaders(r Request) response.Response {
	w := response.BadRequest(errFoo)
	w.WithHeader("lorem-ipsum", "dolor-sit-amet")

	return w
}

func bodyToResponse(r Request) response.Response {
	body := customBodyWithoutValidation{}
	if err := r.Body(&body); err != nil {
		return response.BadRequest(err)
	}

	fmt.Printf("%+v", body)

	return response.Ok(body)
}

func bodyToResponseWithValidation(r Request) response.Response {
	body := customBody{}
	if err := r.Body(&body); err != nil {
		return response.BadRequest(err)
	}

	return response.Ok(body)
}

func redirectHandler(r Request) response.Response {
	return response.MovedPermanently("http://www.onet.pl")
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

type customBody struct {
	A string `json:"a" xml:"a" validate:"iscolor"`
}
type customBodyWithoutValidation struct {
	A string `json:"a" xml:"a"`
}
