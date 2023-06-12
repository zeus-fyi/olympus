package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

const (
	V3_SWAP_EXACT_IN            = 0x00
	V3_SWAP_EXACT_OUT           = 0x01
	PERMIT2_TRANSFER_FROM       = 0x02
	PERMIT2_PERMIT_BATCH        = 0x03
	SWEEP                       = 0x04
	TRANSFER                    = 0x05
	PAY_PORTION                 = 0x06
	V2_SWAP_EXACT_IN            = 0x08
	V2_SWAP_EXACT_OUT           = 0x09
	PERMIT2_PERMIT              = 0x0a
	WRAP_ETH                    = 0x0b
	UNWRAP_WETH                 = 0x0c
	PERMIT2_TRANSFER_FROM_BATCH = 0x0d
	SEAPORT                     = 0x10
	LOOKS_RARE_721              = 0x11
	NFTX                        = 0x12
	CRYPTOPUNKS                 = 0x13
	LOOKS_RARE_1155             = 0x14
	OWNER_CHECK_721             = 0x15
	OWNER_CHECK_1155            = 0x16
	SWEEP_ERC721                = 0x17
	X2Y2_721                    = 0x18
	SUDOSWAP                    = 0x19
	NFT20                       = 0x1a
	X2Y2_1155                   = 0x1b
	FOUNDATION                  = 0x1c
	SWEEP_ERC1155               = 0x1d
)

const (
	Sweep      = "SWEEP"
	PayPortion = "PAY_PORTION"
	Transfer   = "TRANSFER"
	UnwrapWETH = "UNWRAP_WETH"
	WrapETH    = "WRAP_ETH"
)

var (
	UniversalRouterDecoderAbi = MustLoadUniversalRouterDecodingAbi()
	UniversalRouterAbi        = MustLoadUniversalRouterAbi()
)

type UnwrapWETHParams struct {
	Recipient accounts.Address `json:"recipient"`
	AmountMin *big.Int         `json:"amountMin"`
}

func (u *UnwrapWETHParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[UnwrapWETH].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountMin, err := ParseBigInt(args["amountMin"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	u.Recipient = recipient
	u.AmountMin = amountMin
	return nil
}

func (u *UnwrapWETHParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[UnwrapWETH].Inputs.Pack(u.Recipient, u.AmountMin)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

type WrapETHParams struct {
	Recipient accounts.Address `json:"recipient"`
	AmountMin *big.Int         `json:"amountMin"`
}

func (w *WrapETHParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[WrapETH].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountMin, err := ParseBigInt(args["amountMin"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	w.Recipient = recipient
	w.AmountMin = amountMin
	return nil
}

func (w *WrapETHParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[WrapETH].Inputs.Pack(w.Recipient, w.AmountMin)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (t *TransferParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Transfer].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	value, err := ParseBigInt(args["value"])
	if err != nil {
		return err
	}
	token, err := ConvertToAddress(args["token"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	t.Token = token
	t.Recipient = recipient
	t.Value = value
	return nil
}

func (t *TransferParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[Transfer].Inputs.Pack(t.Token, t.Recipient, t.Value)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

type TransferParams struct {
	Token     accounts.Address `json:"token"`
	Recipient accounts.Address `json:"recipient"`
	Value     *big.Int         `json:"value"`
}

type PayPortionParams struct {
	Token     accounts.Address `json:"token"`
	Recipient accounts.Address `json:"recipient"`
	Bips      *big.Int         `json:"bips"`
}

func (p *PayPortionParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[PayPortion].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	value, err := ParseBigInt(args["bips"])
	if err != nil {
		return err
	}
	token, err := ConvertToAddress(args["token"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	p.Token = token
	p.Recipient = recipient
	p.Bips = value
	return nil
}

func (p *PayPortionParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[PayPortion].Inputs.Pack(p.Token, p.Recipient, p.Bips)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}
