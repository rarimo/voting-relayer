package passport_root_update

import (
	"context"
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rarimo/voting-relayer/internal/config"
	"github.com/rarimo/voting-relayer/internal/utils"
	registrationsmtreplicator "github.com/rarimo/voting-relayer/pkg/passport/contracts"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

var (
	proofType, _ = abi.NewType("bytes32[]", "", nil)
	sigType, _   = abi.NewType("bytes", "", nil)
	proofArgs    = abi.Arguments{
		{
			Name: "path",
			Type: proofType,
		},
		{
			Name: "signature",
			Type: sigType,
		},
	}
)

type Core struct {
	log  *logan.Entry
	core rarimocore.QueryClient
	evm  *config.EVM
}

func NewCore(cfg config.Config) *Core {
	return &Core{
		core: rarimocore.NewQueryClient(cfg.Cosmos()),
		log:  cfg.Log().WithField("service", "core"),
		evm:  cfg.EVM(),
	}
}

func (c *Core) GetPassportRootTransferProof(ctx context.Context, operationID string) (*PassportRootTransferDetails, error) {
	proof, err := c.core.OperationProof(ctx, &rarimocore.QueryGetOperationProofRequest{Index: operationID})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the operation proof")
	}

	pathHashes := make([]common.Hash, 0, len(proof.Path))
	for _, p := range proof.Path {
		pathHashes = append(pathHashes, common.HexToHash(p))
	}

	signature := hexutil.MustDecode(proof.Signature)
	signature[64] += 27

	operation, err := c.core.Operation(context.TODO(), &rarimocore.QueryGetOperationRequest{Index: operationID})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the operation")
	}

	transfer, err := pkg.GetPassportRootUpdate(operation.Operation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse operation details")
	}

	result := PassportRootTransferDetails{Operation: transfer}

	result.Proof, err = proofArgs.Pack(pathHashes, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the proof")
	}

	return &result, nil
}

func (c *Core) ProcessPassportStateTransfer(ctx context.Context, details *PassportRootTransferDetails, waitTxConfirm bool) (map[int64]string, error) {
	var commonError error

	res := make(map[int64]string)
	for _, chain := range c.evm.Chains {
		opts := chain.TransactorOpts()

		nonce, err := chain.RPC.PendingNonceAt(context.TODO(), chain.SubmitterAddress)
		if err != nil {
			commonError = errors.Wrap(commonError, errors.Wrap(err, "failed to fetch a nonce").Error())
			continue
		}

		opts.Nonce = big.NewInt(int64(nonce))

		contract, err := registrationsmtreplicator.NewRegistrationSMTReplicatorTransactor(chain.ContractAddress, chain.RPC)
		if err != nil {
			commonError = errors.Wrap(commonError, errors.Wrap(err, "failed to create contract instance").Error())
			continue
		}

		var root [32]byte
		copy(root[:], hexutil.MustDecode(details.Operation.Root))

		timestamp := new(big.Int).SetInt64(details.Operation.Timestamp)

		tx, err := contract.TransitionRoot(opts, root, timestamp, details.Proof)
		if err != nil {
			c.log.Debugf(
				"Tx args: %s, %s, %s",
				root,
				timestamp.String(),
				hexutil.Encode(details.Proof),
			)
			commonError = errors.Wrap(commonError, errors.Wrap(err, "failed to send state transition tx").Error())
			continue
		}

		if waitTxConfirm {
			c.WaitTxConfirmation(ctx, &chain, tx)
		}

		res[chain.ChainID.Int64()] = tx.Hash().Hex()
	}
	return res, commonError
}

func (c *Core) WaitTxConfirmation(ctx context.Context, chain *config.EVMChain, tx *types.Transaction) {
	receipt, err := bind.WaitMined(ctx, chain.RPC, tx)
	if err != nil {
		c.log.WithError(err).Error("failed to wait for state transition tx")
		return
	}

	if receipt.Status == 0 {
		c.log.WithError(err).WithFields(logan.F{
			"receipt": utils.Prettify(receipt),
			"chain":   chain.Name,
		}).Error("failed to wait for state transition tx")
		return
	}

	c.log.
		WithFields(logan.F{
			"tx_id":        tx.Hash(),
			"tx_index":     receipt.TransactionIndex,
			"block_number": receipt.BlockNumber,
			"gas_used":     receipt.GasUsed,
		}).
		Info("evm transaction confirmed")
}
