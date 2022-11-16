package athena_jwt

import (
	"testing"

	"github.com/stretchr/testify/suite"
	aegis_jwt "github.com/zeus-fyi/olympus/pkg/aegis/jwt"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type SetJwtTestSuite struct {
	base.TestSuite
}

func (s *SetJwtTestSuite) TestSetJwtToken() {
	p := structs.Path{DirIn: "./", DirOut: "./", FnIn: ""}

	j := aegis_jwt.NewAegisJWT()
	j.GenerateJwtTokenString()

	err := SetToken(p, j.Token.Raw)
	s.Require().Nil(err)
}

func TestSetJwtTestSuite(t *testing.T) {
	suite.Run(t, new(SetJwtTestSuite))
}
