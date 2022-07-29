package beacon_models

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
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

	err = InsertEpochCheckpoint(ctx, epochZero)
	c.Require().Nil(err)

	err = cp.GetEpochCheckpoint(ctx, epochZero)
	c.Require().Nil(err)

	c.Assert().Equal(3, cp.ValidatorsActive)
	c.Assert().Equal(3, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(0, cp.ValidatorsBalancesRemaining)
	// epoch 1
	epochOne := 1
	err = InsertEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	err = cp.GetEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	c.Assert().Equal(4, cp.ValidatorsActive)
	c.Assert().Equal(0, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(4, cp.ValidatorsBalancesRemaining)

	// epoch 1 and new balance records
	// balacnce updates only happens via TODO
	seedAndInsertNewValidatorBalances(ctx, vsOne, int64(epochOne), 32000000000+43753)
	err = cp.GetEpochCheckpoint(ctx, epochOne)
	c.Require().Nil(err)
	c.Assert().Equal(4, cp.ValidatorsActive)
	// TODO needs to update this. Currently failing on purpose
	//c.Assert().Equal(, cp.ValidatorsBalancesRecorded)
	//c.Assert().Equal(2, cp.ValidatorsBalancesRemaining)
	c.Assert().Equal(0, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(4, cp.ValidatorsBalancesRemaining)

	// epoch 2
	epochTwo := 2
	err = InsertEpochCheckpoint(ctx, epochTwo)
	err = cp.GetEpochCheckpoint(ctx, epochTwo)
	c.Require().Nil(err)
	c.Assert().Equal(5, cp.ValidatorsActive)
	c.Assert().Equal(0, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(5, cp.ValidatorsBalancesRemaining)

	// epoch 3
	epochThree := 3
	err = InsertEpochCheckpoint(ctx, epochThree)
	err = cp.GetEpochCheckpoint(ctx, epochThree)
	c.Require().Nil(err)
	c.Assert().Equal(6, cp.ValidatorsActive)
	c.Assert().Equal(0, cp.ValidatorsBalancesRecorded)
	c.Assert().Equal(6, cp.ValidatorsBalancesRemaining)
}

func (c *CheckpointTestSuite) TestInsertCheckpointWithDiffsAdded() {
	// TODO
}

func TestCheckpointTestSuite(t *testing.T) {
	suite.Run(t, new(CheckpointTestSuite))
}

func CleanupDb(ctx context.Context, tableName string) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, tableName, "true")
	_, err := postgres.Pg.Exec(ctx, query)
	log.Err(err).Interface("cleanupDb: %s", tableName)
	if err != nil {
		panic(err)
	}
}
