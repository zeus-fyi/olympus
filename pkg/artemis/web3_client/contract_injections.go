package web3_client

import artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"

const (
	RawDawgAddr = "0x7623e9dc0da6ff821ddb9ebaba794054e078f8c4"
)

func (w *Web3Client) MustInjectRawDawg() {
	err := w.SetCodeOverride(ctx, RawDawgAddr, artemis_oly_contract_abis.RawdawgByteCode)
	if err != nil {
		panic(err)
	}
	return
}
