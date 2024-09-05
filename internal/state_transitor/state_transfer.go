package ingester

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rarimo/voting-relayer/internal/config"
	"github.com/rarimo/voting-relayer/internal/data"
	passportrootupdate "github.com/rarimo/voting-relayer/internal/state_transitor/core/passport_root_update"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type stateIngester struct {
	log        *logan.Entry
	rarimocore rarimocore.QueryClient
	storage    data.StateQ
	core       *passportrootupdate.Core
}

var _ Processor = &stateIngester{}

func NewPassportRootIngester(cfg config.Config, q data.StateQ) Processor {
	return &stateIngester{
		log:        cfg.Log(),
		rarimocore: rarimocore.NewQueryClient(cfg.Cosmos()),
		core:       passportrootupdate.NewCore(cfg),
		storage:    q,
	}
}

func (s *stateIngester) Query() string {
	return stateQuery
}

func (s *stateIngester) Name() string {
	return "passport-root-ingester"
}

func (s *stateIngester) Catchup(ctx context.Context) error {
	s.log.Info("Starting catchup")
	defer s.log.Info("Catchup finished")

	var nextKey []byte

	for {
		operations, err := s.rarimocore.OperationAll(ctx, &rarimocore.QueryAllOperationRequest{Pagination: &query.PageRequest{Key: nextKey}})
		if err != nil {
			return errors.Wrap(err, "failed to fetch operations")
		}

		for _, op := range operations.Operation {
			if op.Status == rarimocore.OpStatus_SIGNED && op.OperationType == rarimocore.OpType_PASSPORT_ROOT_UPDATE {
				if err = s.Validate(op); err != nil {
					return err
				}
			}
		}

		nextKey = operations.Pagination.NextKey
		if nextKey == nil {
			return nil
		}
	}
}

func (s *stateIngester) Process(
	ctx context.Context,
	confirmationID string,
) error {
	s.log.WithField("confirmation_id", confirmationID).Info("processing a confirmation")

	rawConf, err := s.rarimocore.Confirmation(ctx, &rarimocore.QueryGetConfirmationRequest{Root: confirmationID})
	if err != nil {
		return errors.Wrap(err, "failed to get confirmation", logan.F{
			"confirmation_id": confirmationID,
		})
	}

	for _, index := range rawConf.Confirmation.Indexes {
		operation, err := s.rarimocore.Operation(ctx, &rarimocore.QueryGetOperationRequest{Index: index})
		if err != nil {
			return errors.Wrap(err, "failed to get operation", logan.F{
				"operation_index": operation.Operation.Index,
			})
		}

		if err = s.Validate(operation.Operation); err != nil {
			return err
		}

		proof, err := s.core.GetPassportRootTransferProof(ctx, operation.Operation.Index)
		if err != nil {
			return errors.Wrap(err, "failed to get proof for the operation", logan.F{
				"operation_index": operation.Operation.Index,
			})
		}

		processedOperations, err := s.core.ProcessPassportStateTransfer(ctx, proof, true)
		if err != nil {
			return errors.Wrap(err, "failed to transit proofs", logan.F{
				"operation_index": operation.Operation.Index,
			})
		}

		var commonError error

		for chain, txHash := range processedOperations {
			_, err = s.storage.Insert(
				data.State{
					ChainID:     chain,
					OperationID: operation.Operation.Index,
					TxHash:      txHash,
					Event:       operation.Operation.OperationType.String(),
					Proof:       hexutil.Encode(proof.Proof),
					Root:        proof.Operation.Root,
					BlockHeight: proof.Operation.BlockHeight,
				},
			)
			if err != nil {
				commonError = errors.Wrap(commonError, errors.Wrap(err, "failed to insert state", logan.F{}).Error())
			}
		}

		if commonError != nil {
			return errors.Wrap(commonError, "failed to insert operation info into DB", logan.F{
				"operation_index": operation.Operation.Index,
			})
		}
	}

	return nil
}

func (s *stateIngester) Validate(operation rarimocore.Operation) error {
	if operation.OperationType != rarimocore.OpType_PASSPORT_ROOT_UPDATE {
		return nil
	}

	s.log.WithField("operation_index", operation.Index).Info("Trying to save op")

	if _, err := pkg.GetPassportRootUpdate(operation); err != nil {
		return errors.Wrap(err, "failed to parse passport root transfer", logan.F{
			"operation_index": operation.Index,
		})
	}

	return nil
}
