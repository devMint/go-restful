package restful

import (
	"bytes"
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

const (
	readCollection = "read collection"
	createElement  = "create element"
	readElement    = "read element"
	updateElement  = "update element"
	removeElement  = "remove element"
)

func Test_BasicCRUD_ReadCollection(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/products", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), readCollection)
}

func Test_BasicCRUD_CreateElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/products", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Contains(t, response.Body.String(), createElement)
}

func Test_BasicCRUD_ReadElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/products/1", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), readElement)
}

func Test_BasicCRUD_OverwriteElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/products/1", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), updateElement)
}

func Test_BasicCRUD_UpdateElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PATCH", "/products/1", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), updateElement)
}

func Test_BasicCRUD_DeleteElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/products/1", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func Test_ArticleCRUD_ReadCollection(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/articles", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "{\"id\":1,\"title\":\"lorem-ipsum\"}")
}

func Test_ArticleCRUD_CreateElement_InvalidPayload(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/articles", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "empty body from request")
}

func Test_ArticleCRUD_CreateElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/articles", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "invalid payload")
}

func Test_ArticleCRUD_ReadElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/articles/1", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "{\"id\":1,\"title\":\"lorem-ipsum\"}")
}

func Test_ArticleCRUD_OverwriteElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/articles/1", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), updateElement)
}

func Test_ArticleCRUD_UpdateElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PATCH", "/articles/1", validPayloadArticle())
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), updateElement)
}

func Test_ArticleCRUD_DeleteElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/articles/1", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func createCRUDRouter() Router {
	router := NewRouter(chi.NewMux())
	router.Mount("/products", NewCRUD(crudService{}))
	router.Mount("/articles", NewCRUD(articleService{}))

	return router
}

func validPayloadArticle() *bytes.Buffer {
	body, _ := json.Marshal(map[string]string{"title": "Dolor sit amet"})
	return bytes.NewBuffer(body)
}

type crudService struct{}
type article struct {
	ID    int    `json:"id" xml:"id"`
	Title string `json:"title" xml:"title"`
}

func (c crudService) FindOne(id string) (interface{}, error)  { return readElement, nil }
func (c crudService) Find() (interface{}, error)              { return readCollection, nil }
func (c crudService) Create(interface{}) (interface{}, error) { return createElement, nil }
func (c crudService) Update(interface{}) (interface{}, error) { return updateElement, nil }
func (c crudService) Delete(interface{}) error                { return nil }
func (c crudService) Model() interface{}                      { return article{} }
func (c crudService) CreatePayloadModel() interface{}         { return article{} }
func (c crudService) UpdatePayloadModel() interface{}         { return article{} }

type articleService struct{}
type payloadArticle struct {
	Title string `json:"title" xml:"title"`
}

func (a articleService) FindOne(id string) (interface{}, error) {
	return article{ID: 1, Title: "lorem-ipsum"}, nil
}

func (a articleService) Find() (interface{}, error) {
	return []article{{ID: 1, Title: "lorem-ipsum"}}, nil
}

func (a articleService) Create(data interface{}) (interface{}, error) {
	article, ok := data.(payloadArticle)
	if !ok {
		return nil, response.BadRequest(errors.New("conversion to type 'payloadArticle' not work"))
	}

	return article, nil
}

func (a articleService) Update(data interface{}) (interface{}, error) {
	article, ok := data.(payloadArticle)
	if !ok {
		return nil, response.BadRequest(errors.New("conversion to type 'payloadArticle' not work"))
	}

	return article, nil
}

func (a articleService) Delete(interface{}) error { return nil }

func (a articleService) Model() interface{} { return &article{} }

func (c articleService) CreatePayloadModel() interface{} { return payloadArticle{} }

func (c articleService) UpdatePayloadModel() interface{} { return payloadArticle{} }

type exampleService struct{}

func (e exampleService) FindOne() (article, error) {
	return article{ID: 1, Title: "lorem-ipsum"}, nil
}
func (e exampleService) Find() ([]article, error) {
	return []article{{ID: 1, Title: "lorem-ipsum"}}, nil
}
func (e exampleService) Create(article article, payload payloadArticle) (article, error) {
	article.Title = payload.Title
	return article, nil
}

func ArticleFindOne(r request.Request) response.Response {
	article, ok := r.Context().Value("article").(article)
	if !ok {
		return response.NotFound(errors.New("entity not found"))
	}

	return response.Ok(article)
}
