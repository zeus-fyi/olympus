package demo

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func DemoPath() filepaths.Path {
	var demoPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./demo",
		DirOut:      "./demo_out/gzip",
		FnIn:        "demo",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return demoPath
}

func DemoReadChartThenWritePath() filepaths.Path {
	var demoPath = filepaths.Path{
		PackageName: "",
		DirIn:       "./demo",
		DirOut:      "./demo_out/read_chart",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	return demoPath
}
