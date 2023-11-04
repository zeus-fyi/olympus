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
		blockNumber := c.GetPendingMainnetSlotNum()
		slotNumberOffset := c.GetSecsSinceLastMainnetSlot()
		fmt.Println(slotNumber, blockNumber, slotNumberOffset)
	}
}

func (s *ChronosTestSuite) TestConvertToTime() {
	c := Chronos{}
	nt := c.ConvertUnixTimeStampToDate(1699137400414158691)
	s.Require().NotEmpty(nt)
	fmt.Println(nt.String())
}
func TestChronosTestSuite(t *testing.T) {
	suite.Run(t, new(ChronosTestSuite))
}
