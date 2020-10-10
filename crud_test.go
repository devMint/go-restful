package restful

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

const (
	readCollection   = "read collection"
	createElement    = "create element"
	readElement      = "read element"
	updateElement    = "update element"
	overwriteElement = "overwrite element"
	removeElement    = "remove element"
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
	request, _ := http.NewRequest("POST", "/products", nil)
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
	request, _ := http.NewRequest("PUT", "/products/1", nil)
	createCRUDRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), overwriteElement)
}

func Test_BasicCRUD_UpdateElement(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PATCH", "/products/1", nil)
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

func createCRUDRouter() Router {
	router := NewRouter(chi.NewMux())
	router.Mount("/products", NewCRUD(crudService{}))

	return router
}

type crudService struct{}

func (c crudService) FindOne(id string) (interface{}, error)  { return readCollection, nil }
func (c crudService) Find() (interface{}, error)              { return readElement, nil }
func (c crudService) Create(interface{}) (interface{}, error) { return createElement, nil }
func (c crudService) Update(interface{}) (interface{}, error) { return updateElement, nil }
func (c crudService) Delete(interface{}) error                { return nil }
