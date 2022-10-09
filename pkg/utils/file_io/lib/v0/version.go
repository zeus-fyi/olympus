package v0

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
)

type Lib struct {
	readers.ReaderLib
	paths.PathLib
}
