package handlers

import (
	"encoding/hex"
	"github.com/rarimo/voting-relayer/internal/service/requests"
	"github.com/rarimo/voting-relayer/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func Operation(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetOperationIDRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	operation, err := StateQ(r).FilterByOperationId(req.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get blob from DB")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if operation == nil {
		Log(r).WithError(err).Error("operation not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, resources.Relation{
		Data: &resources.Key{
			ID:   hex.EncodeToString(operation.TxHash[:]),
			Type: resources.TRANSACTION,
		},
	})
}
