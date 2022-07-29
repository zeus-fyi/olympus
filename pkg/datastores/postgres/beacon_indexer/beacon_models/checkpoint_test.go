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

func (c *CheckpointTestSuite) TestInsertCheckpoint() {
	ctx := context.Background()
	cp := ValidatorsEpochCheckpoint{
		Epoch: 0,
	}

	numValidators := 3
	beforeInsertCount, err := SelectCountValidatorActive(ctx, cp.Epoch)
	c.Assert().Equal(0, beforeInsertCount)

	createAndInsertFakeValidatorsByStatus(numValidators, "active_ongoing_genesis")
	firstInsertCount, err := SelectCountValidatorActive(ctx, cp.Epoch)
	c.Assert().Equal(numValidators, firstInsertCount)

	err = InsertEpochCheckpoint(ctx, cp.Epoch)
	c.Require().Nil(err)

	err = cp.GetEpochCheckpoint(ctx, cp.Epoch)
	c.Require().Nil(err)

	c.Assert().Equal(numValidators, cp.ValidatorsActive)
	c.Assert().Equal(3, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(0, cp.ValidatorsBalancesRemaining)
}

func TestCheckpointTestSuite(t *testing.T) {
	suite.Run(t, new(CheckpointTestSuite))
}
