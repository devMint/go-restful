Small wrapper for [go-chi/chi](https://github.com/go-chi/chi) to make HTTP handlers more return-like. It doesn't add new functionality to router, it just allows to easier handle responses.

**This is work in progress package!**

## Summary

* `restful.NewRouter()` requires an instance of `chi.Router`. This allows to use existing router without breaking the whole codebase.
* existing `chiRouter` remains existing routes and middlewares so it allows you to use `http.Handler` and `http.HandlerFunc` with chi router
* restful's routes use new `request.ContextAction` and `request.RestfulHandler` definitions:

```go
type RestfulHandler func(Request) response.Response
type ContextHandler func(Request) (context.Context, response.Response)
```

## Examples

```go
package main

import (
	"net/http"

	"github.com/devMint/go-restful"
	"github.com/devMint/go-restful/request"
	"github.com/devMint/go-restful/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// This is the old way to define router and its handlers.
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)
	chiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	restfulRouter := restful.NewRouter(chiRouter)
	restfulRouter.Get("/", func(r request.Request) response.Response {
		return response.Ok("welcome")
	})

	// http.ListenAndServe(":3000", chiRouter)
	http.ListenAndServe(":3000", restfulRouter)
}
```