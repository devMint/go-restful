package restful

import (
	"context"
	"errors"

	"github.com/devMint/go-restful/request"
	"github.com/devMint/go-restful/response"
	"github.com/go-chi/chi"
)

type CRUD interface {
	FindOne(id string) (interface{}, error)
	Find() (interface{}, error)
	Create(interface{}) (interface{}, error)
	Update(interface{}) (interface{}, error)
	Delete(interface{}) error
}

func NewCRUD(service CRUD) Router {
	endpoints := crudRoutes{c: service}

	router := NewRouter(chi.NewMux())
	router.Get("/", endpoints.FindAll)
	router.Post("/", endpoints.Create)
	router.Route("/{id}", func(r Router) {
		r.Use(endpoints.WithElement)
		r.Get("/", endpoints.FindOne)
		r.Put("/", endpoints.Update)
		r.Patch("/", endpoints.Update)
		r.Delete("/", endpoints.Delete)
	})

	return router
}

type crudRoutes struct {
	c CRUD
}

func (crud crudRoutes) FindAll(r request.Request) response.Response {
	data, err := crud.c.Find()
	if err != nil {
		return response.BadRequest(err)
	}

	return response.Ok(data)
}

func (crud crudRoutes) FindOne(r request.Request) response.Response {
	data := r.Context().Value("found_element_ctx_crud")
	if data == nil {
		return response.NotFound(errors.New("Entity not found"))
	}

	return response.Ok(data)
}

func (crud crudRoutes) Create(r request.Request) response.Response {
	data := r.Context().Value("found_element_ctx_crud")
	if data == nil {
		return response.NotFound(errors.New("Entity not found"))
	}

	created, err := crud.c.Create(data)
	if err != nil {
		return response.BadRequest(err)
	}

	return response.Created(created)
}

func (crud crudRoutes) Update(r request.Request) response.Response {
	data := r.Context().Value("found_element_ctx_crud")
	if data == nil {
		return response.NotFound(errors.New("Entity not found"))
	}

	created, err := crud.c.Update(data)
	if err != nil {
		return response.BadRequest(err)
	}

	return response.Ok(created)
}

func (crud crudRoutes) Delete(r request.Request) response.Response {
	data := r.Context().Value("found_element_ctx_crud")
	if data == nil {
		return response.NotFound(errors.New("Entity not found"))
	}

	if err := crud.c.Delete(data); err != nil {
		return response.BadRequest(err)
	}

	return response.NoContent()
}

func (crud crudRoutes) WithElement(req request.Request) (context.Context, response.Response) {
	id := req.Param("id")
	data, err := crud.c.FindOne(id)
	if err != nil {
		return req.Context(), response.NotFound(err)
	}

	return context.WithValue(req.Context(), "found_element_ctx_crud", data), nil
}