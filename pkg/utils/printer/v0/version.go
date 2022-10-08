package v0

import (
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

type Lib struct {
	structs.Path
	Log logging.Logger
}
