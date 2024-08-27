package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rarimo/voting-relayer/internal/config"
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
		return fmt.Errorf("failed to suggest gas price: %w", err)
	}

	txd.Gas, err = relayerConfig.RPC.EstimateGas(ctx, ethereum.CallMsg{
		From:     crypto.PubkeyToAddress(relayerConfig.PrivateKey.PublicKey),
		To:       receiver,
		GasPrice: txd.GasPrice,
		Data:     txd.DataBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	return nil
}

func SendTx(ctx context.Context, txd *TxData, receiver *common.Address, relayerConfig *config.RelayerConfig) (tx *types.Transaction, err error) {
	tx, err = SignTx(txd, receiver, relayerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to sign new tx: %w", err)
	}

	if err = relayerConfig.RPC.SendTransaction(ctx, tx); err != nil {
		if strings.Contains(err.Error(), "nonce") {
			if err = relayerConfig.ResetNonce(relayerConfig.RPC); err != nil {
				return nil, fmt.Errorf("failed to reset nonce: %w", err)
			}

			tx, err = SignTx(txd, receiver, relayerConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to sign new tx: %w", err)
			}

			if err := relayerConfig.RPC.SendTransaction(ctx, tx); err != nil {
				return nil, fmt.Errorf("failed to send transaction: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to send transaction: %w", err)
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
