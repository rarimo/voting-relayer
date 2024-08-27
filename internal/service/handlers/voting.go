package handlers

import (
	"github.com/rarimo/voting-relayer/internal/utils"
	"math/big"
	"net/http"
	"strings"

	"github.com/rarimo/voting-relayer/internal/service/proposalsstate"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/voting-relayer/internal/service/requests"
	"github.com/rarimo/voting-relayer/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Voting(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVotingRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to get request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var (
		destination = req.Data.Attributes.Destination
		calldata    = req.Data.Attributes.TxData
		proposalID  = req.Data.Attributes.ProposalId
	)

	log := Log(r).WithFields(logan.F{
		"user-agent":  r.Header.Get("User-Agent"),
		"calldata":    calldata,
		"destination": destination,
		"proposal_id": proposalID,
	})
	log.Debug("voting request")

	// destination is valid hex address because of request validation
	votingAddress := common.HexToAddress(destination)

	var txd utils.TxData
	txd.DataBytes, err = hexutil.Decode(calldata)
	if err != nil {
		log.WithError(err).Error("Failed to decode data")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proposalBigID := big.NewInt(proposalID)

	session, err := proposalsstate.NewProposalsStateCaller(RelayerConfig(r).Address, RelayerConfig(r).RPC)

	if err != nil {
		log.WithError(err).Error("Failed to get proposal state caller")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	proposalConfig, err := session.GetProposalConfig(nil, proposalBigID)

	if err != nil {
		log.WithError(err).Error("Failed to get proposal config")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !utils.IsAddressInWhitelist(votingAddress, proposalConfig.VotingWhitelist) {
		log.Error("Address not in voting whitelist")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	defer RelayerConfig(r).UnlockNonce()
	RelayerConfig(r).LockNonce()

	err = utils.ConfGas(r.Context(), &txd, &votingAddress, RelayerConfig(r))

	if err != nil {
		log.WithError(err).Error("Failed to configure gas and gasPrice")
		// `errors.Is` is not working for rpc errors, they passed as a string without additional wrapping
		// because of this we operate with raw strings
		if strings.Contains(err.Error(), vm.ErrExecutionReverted.Error()) {
			errParts := strings.Split(err.Error(), ":")
			contractName := strings.TrimSpace(errParts[len(errParts)-2])
			errMsg := errors.New(strings.TrimSpace(errParts[len(errParts)-1]))
			ape.RenderErr(w, problems.BadRequest(validation.Errors{contractName: errMsg}.Filter())...)
			return
		}
		ape.RenderErr(w, problems.InternalError())
		return
	}

	tx, err := utils.SendTx(r.Context(), &txd, &votingAddress, RelayerConfig(r))
	if err != nil {
		log.WithError(err).Error("failed to send tx")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	RelayerConfig(r).IncrementNonce()

	ape.Render(w, resources.Relation{
		Data: &resources.Key{
			ID:   tx.Hash().String(),
			Type: resources.TRANSACTION,
		},
	})
}
