package handlers

import (
	"github.com/rarimo/voting-relayer/internal/data"
	"github.com/rarimo/voting-relayer/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func GetLastOperation(w http.ResponseWriter, r *http.Request) {
	operation, err := StateQ(r).SortByBlockHeight(data.DESC).Get()
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
				DestinationChain: strconv.FormatInt(operation.ChainId, 10),
				OperationId:      operation.OperationId,
				Proof:            operation.Proof,
				TxHash:           operation.TxHash,
			},
		},
	})
}
