package requests

import (
	"net/http"

	"github.com/go-chi/chi"
)

type GetOperationByRootRequest struct {
	Root string
}

func NewGetOperationByRootRequest(r *http.Request) (GetOperationByRootRequest, error) {
	request := GetOperationByRootRequest{}
	request.Root = chi.URLParam(r, "root")

	return request, nil
}
