package chronos

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ChronosTestSuite struct {
	suite.Suite
}

func (s *ChronosTestSuite) TestLib0() {
	c := Chronos{}
	s.Require().NotEmpty(c.UnixTimeStampNow())
}

func (s *ChronosTestSuite) TestEthUtils() {

	c := Chronos{}
	for i := 0; i < 100; i++ {
		time.Sleep(200 * time.Millisecond)
		slotNumber := c.GetPendingMainnetSlotNum()
		blockNumber := c.GetPendingMainnetBlockNumber()
		slotNumberOffset := c.GetSecsSinceLastMainnetSlot()
		fmt.Println(slotNumber, blockNumber, slotNumberOffset)
	}
}

func TestChronosTestSuite(t *testing.T) {
	suite.Run(t, new(ChronosTestSuite))
}
