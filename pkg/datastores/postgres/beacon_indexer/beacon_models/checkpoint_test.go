package beacon_models

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type CheckpointTestSuite struct {
	test_suites.PGTestSuite
	P *pgxpool.Pool
}

func (c *CheckpointTestSuite) InsertCheckpoint() {
	ctx := context.Background()

	validators := 10
	vs := createFakeValidatorsByStatus(validators, "active_ongoing")
	err := vs.InsertValidatorsPendingQueue(ctx)
	c.Require().Nil(err)

	cp := EpochCheckpoint{
		Epoch: 0,
	}
	err = InsertEpochCheckpoint(ctx, cp.Epoch)
	c.Require().Nil(err)

	err = cp.GetEpochCheckpoint(ctx, cp.Epoch)
	c.Require().Nil(err)
	c.Assert().Equal(validators, cp.ValidatorsActive)
	c.Assert().Equal(0, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(validators, cp.ValidatorsBalancesRemaining)
}

func TestCheckpointTestSuite(t *testing.T) {
	suite.Run(t, new(CheckpointTestSuite))
}
