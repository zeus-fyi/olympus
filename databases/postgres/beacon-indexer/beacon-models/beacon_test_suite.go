package beacon_models

import (
	"math/rand"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/test_utils/test_suites"
)

type BeaconBaseTestSuite struct {
	test_suites.PGTestSuite
	P *pgxpool.Pool
}

func createFakeValidators(validatorNum int) Validators {
	var vs Validators
	for i := 0; i < validatorNum; i++ {
		var val Validator
		val.Index = rand.Int63n(100000)
		val.Pubkey = "0x" + random.String(96)
		vs.Validators = append(vs.Validators, val)
	}
	return vs
}

func TestBeaconBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconBaseTestSuite))
}
