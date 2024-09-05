package handlers

import (
	"net/http"
	"strconv"

	"github.com/rarimo/voting-relayer/internal/service/requests"
	"github.com/rarimo/voting-relayer/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetOperationByRoot(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetOperationByRootRequest(r)
	if err != nil {
		Log(r).WithError(err).Info("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	operation, err := StateQ(r).FilterByRoot(req.Root).Get()
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

	ape.Render(w, resources.OperationResponse{
		Data: resources.Operation{
			Attributes: resources.OperationAttributes{
				BlockHeight:      operation.BlockHeight,
				DestinationChain: strconv.FormatInt(operation.ChainID, 10),
				OperationID:      operation.OperationID,
				Proof:            operation.Proof,
				TxHash:           operation.TxHash,
			},
		},
	})
}
