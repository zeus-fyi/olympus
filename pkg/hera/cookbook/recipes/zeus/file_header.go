package zeus

import (
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

var autogenPath = structs.Path{
	PackageName: "autogen_structs",
	DirIn:       "",
	DirOut:      "postgres/apps/zeus/structs/autogen_preview",
	Fn:          "",
	Env:         "",
}

var baseFw = primitives.FileWrapper{}
