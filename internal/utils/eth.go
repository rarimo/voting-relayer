package utils

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rarimo/voting-relayer/internal/config"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
	"strings"
)

type TxData struct {
	DataBytes []byte
	GasPrice  *big.Int
	Gas       uint64
}

func IsAddressInWhitelist(votingAddress common.Address, whitelist []common.Address) bool {
	votingAddressBytes := votingAddress.Bytes()
	for _, addr := range whitelist {
		if bytes.Equal(votingAddressBytes, addr.Bytes()) {
			return true
		}
	}
	return false
}

func ConfGas(ctx context.Context, txd *TxData, receiver *common.Address, relayerConfig *config.RelayerConfig) (err error) {
	txd.GasPrice, err = relayerConfig.RPC.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to suggest gas price")
	}

	txd.Gas, err = relayerConfig.RPC.EstimateGas(ctx, ethereum.CallMsg{
		From:     crypto.PubkeyToAddress(relayerConfig.PrivateKey.PublicKey),
		To:       receiver,
		GasPrice: txd.GasPrice,
		Data:     txd.DataBytes,
	})
	if err != nil {
		return errors.Wrap(err, "failed to estimate gas")
	}

	return nil
}

func SendTx(ctx context.Context, txd *TxData, receiver *common.Address, relayerConfig *config.RelayerConfig) (tx *types.Transaction, err error) {
	tx, err = SignTx(txd, receiver, relayerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign new tx")
	}

	if err = relayerConfig.RPC.SendTransaction(ctx, tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			if err = relayerConfig.ResetNonce(relayerConfig.RPC); err != nil {
				return nil, errors.Wrap(err, "failed to reset nonce")
			}

			tx, err = SignTx(txd, receiver, relayerConfig)
			if err != nil {
				return nil, errors.Wrap(err, "failed to sign new tx")
			}

			if err := relayerConfig.RPC.SendTransaction(ctx, tx); err != nil {
				return nil, errors.Wrap(err, "failed to send transaction")
			}
		} else {
			return nil, errors.Wrap(err, "failed to send transaction")
		}
	}

	return tx, nil
}

func SignTx(txd *TxData, receiver *common.Address, relayerConfig *config.RelayerConfig) (tx *types.Transaction, err error) {
	tx, err = types.SignNewTx(
		relayerConfig.PrivateKey,
		types.NewCancunSigner(relayerConfig.ChainID),
		&types.LegacyTx{
			Nonce:    relayerConfig.Nonce(),
			Gas:      txd.Gas,
			GasPrice: txd.GasPrice,
			To:       receiver,
			Data:     txd.DataBytes,
		},
	)
	return tx, err
}
