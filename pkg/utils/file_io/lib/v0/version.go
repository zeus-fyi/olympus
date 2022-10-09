package v0

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
)

type Lib struct {
	structs.Path
	Log logging.Logger
}
