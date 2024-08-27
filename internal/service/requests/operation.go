package requests

import (
	"github.com/go-chi/chi"
	"net/http"
)

type GetOperationByIDRequest struct {
	ID string
}

func NewGetOperationIDRequest(r *http.Request) (GetOperationByIDRequest, error) {
	request := GetOperationByIDRequest{}
	request.ID = chi.URLParam(r, "id")
	return request, nil
}
