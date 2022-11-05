package zeus

import (
	"bytes"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

func UnGzipK8sChart(in *bytes.Buffer) (chart_workload.NativeK8s, error) {
	yr := transformations.YamlFileIO{}
	p := structs.Path{DirIn: "/tmp", DirOut: "/tmp", Fn: "chart.tar.gz"}
	m := memfs.NewMemFs()
	err := m.MakeFile(&p, in.Bytes())
	if err != nil {
		return yr.NativeK8s, err
	}
	p.DirOut = "/chart"
	comp := compression.NewCompression()
	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
	if err != nil {
		return yr.NativeK8s, err
	}

	p.DirIn = "/chart"
	err = yr.ReadK8sWorkloadInMemFsDir(p, m)
	return yr.NativeK8s, err
}
