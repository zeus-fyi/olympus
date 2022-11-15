package aegis_jwt

import (
	"github.com/golang-jwt/jwt/v4"
	aegis_crypto "github.com/zeus-fyi/olympus/pkg/aegis/crypto"
)

type AresJWT struct {
	jwt.Token
}

func NewAresJWT() AresJWT {
	return AresJWT{
		Token: jwt.Token{},
	}
}

func (j *AresJWT) GenerateJwtTokenString() {
	h := aegis_crypto.Hex(32)
	j.Raw = h
}
