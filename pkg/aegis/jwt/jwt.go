package aegis_jwt

import (
	"github.com/golang-jwt/jwt/v4"
	aegis_crypto "github.com/zeus-fyi/olympus/pkg/aegis/crypto"
)

type AegisJWT struct {
	jwt.Token
}

func NewAegisJWT() AegisJWT {
	return AegisJWT{
		Token: jwt.Token{},
	}
}

func (j *AegisJWT) GenerateJwtTokenString() {
	h := aegis_crypto.Hex(32)
	j.Raw = h
}
