package writers

import (
	v0 "github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/file_management"
	"github.com/zeus-fyi/olympus/pkg/utils/logging"
)

type WriterLib struct {
	v0.Lib
	Log logging.Logger
	file_management.FileManagerLib
}
