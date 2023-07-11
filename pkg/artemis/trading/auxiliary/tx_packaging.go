package artemis_trading_auxiliary

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/core/types"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
)

func (a *AuxiliaryTradingUtils) packagePermit2Tx(ctx context.Context, signedTx *types.Transaction, permit2 artemis_autogen_bases.Permit2Tx, userNonceOffset int) (artemis_eth_txs.EthTx, error) {
	nonce, err := a.getNonce(ctx)
	if err != nil {
		return artemis_eth_txs.EthTx{}, err
	}
	mevTx := artemis_eth_txs.EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: permit2.ProtocolNetworkID,
			TxHash:            signedTx.Hash().String(),
			Nonce:             int(nonce) + userNonceOffset,
			From:              a.U.Web3Client.Address().String(),
			Type:              "0x02",
		},
		EthTxGas: artemis_autogen_bases.EthTxGas{
			TxHash: signedTx.Hash().String(),
			GasPrice: sql.NullInt64{
				Valid: false,
			},
			GasLimit: sql.NullInt64{
				Int64: int64(signedTx.Gas()),
				Valid: true,
			},
			GasTipCap: sql.NullInt64{
				Int64: signedTx.GasTipCap().Int64(),
				Valid: true,
			},
			GasFeeCap: sql.NullInt64{
				Int64: signedTx.GasFeeCap().Int64(),
				Valid: true,
			},
		},
		Permit2Tx: artemis_eth_txs.Permit2Tx{
			Permit2Tx: artemis_autogen_bases.Permit2Tx{
				Nonce:             permit2.Nonce,
				Owner:             permit2.Owner,
				Deadline:          permit2.Deadline,
				Token:             permit2.Token,
				ProtocolNetworkID: permit2.ProtocolNetworkID,
			},
		},
	}
	return mevTx, nil
}

func (a *AuxiliaryTradingUtils) packageRegularTx(ctx context.Context, signedTx *types.Transaction, userNonceOffset int) (artemis_eth_txs.EthTx, error) {
	pi := signedTx.ChainId()
	ethGas := artemis_autogen_bases.EthTxGas{
		TxHash: signedTx.Hash().String(),
		GasPrice: sql.NullInt64{
			Valid: false,
		},
		GasLimit: sql.NullInt64{
			Int64: int64(signedTx.Gas()),
			Valid: true,
		},
		GasTipCap: sql.NullInt64{
			Int64: signedTx.GasTipCap().Int64(),
			Valid: true,
		},
		GasFeeCap: sql.NullInt64{
			Int64: signedTx.GasFeeCap().Int64(),
			Valid: true,
		},
	}

	signerType := int(signedTx.Type())
	typeEnum := "0x02"
	if signerType == 1 {
		typeEnum = "0x01"
		ethGas.GasPrice = sql.NullInt64{
			Int64: signedTx.GasPrice().Int64(),
			Valid: true,
		}
		ethGas.GasFeeCap = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
		ethGas.GasTipCap = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
	mevTx := artemis_eth_txs.EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: int(pi.Int64()),
			TxHash:            signedTx.Hash().String(),
			Nonce:             int(signedTx.Nonce()),
			From:              a.U.Web3Client.Address().String(),
			Type:              typeEnum,
		},
		EthTxGas: artemis_autogen_bases.EthTxGas{
			TxHash: signedTx.Hash().String(),
			GasPrice: sql.NullInt64{
				Valid: false,
			},
			GasLimit: sql.NullInt64{
				Int64: int64(signedTx.Gas()),
				Valid: true,
			},
			GasTipCap: sql.NullInt64{
				Int64: signedTx.GasTipCap().Int64(),
				Valid: true,
			},
			GasFeeCap: sql.NullInt64{
				Int64: signedTx.GasFeeCap().Int64(),
				Valid: true,
			},
		},
	}
	return mevTx, nil
}
