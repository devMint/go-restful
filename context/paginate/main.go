package paginate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devMint/go-restful/request"
	"github.com/devMint/go-restful/response"
)

const (
	PaginateTake = "take"
	PaginateSkip = "skip"
)

var (
	takeGreaterThan0 = errors.New("param 'take' should be greater than 0")
	skipGreaterThan0 = errors.New("param 'skip' should be greater than 0")
)

func Paginate(defaultTake int, defaultSkip int) func(http.Handler) http.Handler {
	return request.HandleContext(func(req request.Request) (context.Context, response.Response) {
		take, skip := req.Query("take", fmt.Sprint(defaultTake)), req.Query("skip", fmt.Sprint(defaultSkip))
		takeNum, err := strconv.Atoi(take)
		if err != nil {
			return req.Context(), response.BadRequest(err)
		}
		skipNum, err := strconv.Atoi(skip)
		if err != nil {
			return req.Context(), response.BadRequest(err)
		}

		if takeNum <= 0 {
			return req.Context(), response.BadRequest(takeGreaterThan0)
		}
		if skipNum < 0 {
			return req.Context(), response.BadRequest(skipGreaterThan0)
		}

		ctxTake := context.WithValue(req.Context(), PaginateTake, takeNum)
		ctxSkip := context.WithValue(ctxTake, PaginateSkip, skipNum)

		return ctxSkip, nil
	})
}
