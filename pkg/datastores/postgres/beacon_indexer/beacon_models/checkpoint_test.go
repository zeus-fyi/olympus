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

var tablesToCleanup = []string{"validator_balances_at_epoch", "validators_epoch_checkpoint", "validators"}

func (c *CheckpointTestSuite) TestInsertCheckpoint() {
	ctx := context.Background()
	c.CleanupDb(ctx, tablesToCleanup)

	// epoch 0
	epochZero := 0
	cp := ValidatorsEpochCheckpoint{
		Epoch: epochZero,
	}
	beforeInsertCount, err := SelectCountValidatorActive(ctx, epochZero)
	c.Assert().Equal(0, beforeInsertCount)

	// epoch 0 - 3 validators active
	vsOne := createAndInsertFakeValidatorsByStatus(1, "active_ongoing_genesis")
	_ = createAndInsertFakeValidatorsByStatus(2, "active_ongoing_genesis")
	// epoch 1 - 4 validators_active
	createAndInsertFakeValidatorsByStatus(1, "active_ongoing_at_epoch_one")
	// epoch 2 - 5 validators_active
	createAndInsertFakeValidatorsByStatus(1, "active_ongoing_at_epoch_two")
	// epoch 3 - 6 validators_active
	createAndInsertFakeValidatorsByStatus(1, "active_ongoing_at_epoch_three")

	firstInsertCount, err := SelectCountValidatorActive(ctx, epochZero)
	c.Assert().Equal(3, firstInsertCount)

	// epoch 0
	_, err = InsertEpochCheckpoint(ctx, epochZero)
	c.Require().Nil(err)
	c.assertCheckpointValues(ctx, epochZero, 3, 3, 0)

	// epoch 1
	epochOne := 1
	_, err = InsertEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	c.assertCheckpointValues(ctx, epochOne, 4, 0, 4)

	// epoch 1 and new balance records
	// TODO needs to update this, then check balances
	seedAndInsertNewValidatorBalances(ctx, vsOne, int64(epochOne), 32000000000+43753)
	err = cp.GetEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	c.assertCheckpointValues(ctx, epochOne, 4, 0, 4)

	// epoch 2
	epochTwo := 2
	_, err = InsertEpochCheckpoint(ctx, epochTwo)
	c.assertCheckpointValues(ctx, epochTwo, 5, 0, 5)

	// epoch 3
	epochThree := 3
	_, err = InsertEpochCheckpoint(ctx, epochThree)
	c.assertCheckpointValues(ctx, epochThree, 6, 0, 6)

	var fetchCheckpoint ValidatorsEpochCheckpoint
	err = fetchCheckpoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	c.Require().Nil(err)
	c.Assert().Equal(1, fetchCheckpoint.Epoch)
	c.assertCheckpointValues(ctx, epochOne, 4, 0, 4)

	err = UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx, 1)
	c.Require().Nil(err)

	var updatedCheckpoint ValidatorsEpochCheckpoint
	err = updatedCheckpoint.GetEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	// since trigger adds the first balance entry, this should be counted
	c.assertCheckpointValues(ctx, epochOne, 4, 2, 2)
}

func (c *CheckpointTestSuite) assertCheckpointValues(ctx context.Context, epoch, expValsActive, expValBalancesRecorded, expValBalancesRemaining int) {
	cp := ValidatorsEpochCheckpoint{}
	err := cp.GetEpochCheckpoint(ctx, epoch)
	c.Require().Nil(err)
	c.Assert().Equal(expValsActive, cp.ValidatorsActive)
	c.Assert().Equal(expValBalancesRecorded, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(expValBalancesRemaining, cp.ValidatorsBalancesRemaining)
}

func (c *CheckpointTestSuite) TestInsertCheckpointWithDiffsAdded() {
	// TODO
}

func TestCheckpointTestSuite(t *testing.T) {
	suite.Run(t, new(CheckpointTestSuite))
}
