package aegis_crypto

import (
	"github.com/ethereum/go-ethereum/crypto"
	zeus_ecdsa "github.com/zeus-fyi/zeus/pkg/aegis/crypto/ecdsa"
	aegis_random "github.com/zeus-fyi/zeus/pkg/aegis/crypto/random"
)

func GenAddresses(count, workers int) (zeus_ecdsa.AddressGenerator, error) {
	mnemonic, err := aegis_random.GenerateMnemonic()
	pw := crypto.Keccak256Hash([]byte(mnemonic)).Hex()
	ag, err := zeus_ecdsa.GenerateZeroPrefixAddresses(mnemonic, pw, count, workers)
	if err != nil {
		return ag, err
	}
	return ag, err
}
