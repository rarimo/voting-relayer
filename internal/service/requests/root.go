package requests

import (
	"github.com/go-chi/chi"
	"net/http"
)

type GetOperationByRootRequest struct {
	Root string
}

func NewGetOperationByRootRequest(r *http.Request) (GetOperationByRootRequest, error) {
	request := GetOperationByRootRequest{}
	request.Root = chi.URLParam(r, "root")

	return request, nil
}
