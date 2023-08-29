package iris_usage_meters

import (
	"bytes"
	"io"
	"sync/atomic"
	"time"
)

const (
	ZeusUnitsPerRequest  = 4.0
	ZeusUnitsPerResponse = 4.0
	ZeusUnitsPerKB       = 1.0
	BytesPerKB           = 1024.0
)

type PayloadSizeMeter struct {
	R     io.Reader
	Size  int64
	Plan  string
	Month string
}

type UsageMeter struct {
	RateLimit    float64 `json:"rateLimit"`
	MonthlyUsage float64 `json:"monthlyUsage"`
}

func (u *UsageMeter) IsRateLimited(rateLimit, monthlyLimit float64) (bool, bool) {
	return u.RateLimit > rateLimit, u.MonthlyUsage > monthlyLimit
}

func NewPayloadSizeMeter(bodyBytes []byte) *PayloadSizeMeter {
	month := time.Now().UTC().Month()
	pm := &PayloadSizeMeter{
		R:     bytes.NewReader(bodyBytes),
		Month: month.String(),
	}
	pm.Add(int64(len(bodyBytes)))
	return pm
}

func (cr *PayloadSizeMeter) Read(p []byte) (n int, err error) {
	n, err = cr.R.Read(p)
	if n > 0 {
		atomic.AddInt64(&cr.Size, int64(n))
	}
	return
}

func (cr *PayloadSizeMeter) N() int64 {
	return atomic.LoadInt64(&cr.Size)
}

func (cr *PayloadSizeMeter) Add(n int64) {
	atomic.AddInt64(&cr.Size, n)
}

func (cr *PayloadSizeMeter) SizeInKB() float64 {
	sizeInBytes := atomic.LoadInt64(&cr.Size)
	return float64(sizeInBytes) / BytesPerKB
}

func (cr *PayloadSizeMeter) ZeusRequestComputeUnitsConsumed() float64 {
	sizeInKB := cr.SizeInKB()

	// Add base compute units for the request
	computeUnits := ZeusUnitsPerRequest

	// If the payload size is greater than 1KB, add compute units based on the payload size
	if sizeInKB <= 1 {
		// Subtract 1 because the first KB is already included in the base compute units
		computeUnits += ZeusUnitsPerKB
	} else {
		// If the payload size is less than or equal to 1KB, add 1 Zeus compute unit
		computeUnits += sizeInKB
	}

	return computeUnits
}

func (cr *PayloadSizeMeter) ZeusResponseComputeUnitsConsumed() float64 {
	sizeInKB := cr.SizeInKB()

	// Add base compute units for the request
	computeUnits := ZeusUnitsPerResponse

	// If the payload size is greater than 1KB, add compute units based on the payload size
	if sizeInKB <= 1 {
		// Subtract 1 because the first KB is already included in the base compute units
		computeUnits += ZeusUnitsPerKB
	} else {
		// If the payload size is less than or equal to 1KB, add 1 Zeus compute unit
		computeUnits += sizeInKB
	}

	return computeUnits
}

func (cr *PayloadSizeMeter) Reset() {
	atomic.StoreInt64(&cr.Size, 0)
}
