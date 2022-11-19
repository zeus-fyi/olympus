package zeus

import (
	"bytes"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

func UnGzipK8sChart(in *bytes.Buffer) (chart_workload.TopologyBaseInfraWorkload, error) {
	yr := transformations.YamlFileIO{}
	p := filepaths.Path{DirIn: "/tmp", DirOut: "/tmp", FnIn: "chart.tar.gz"}
	m := memfs.NewMemFs()
	err := m.MakeFileIn(&p, in.Bytes())
	if err != nil {
		return yr.TopologyBaseInfraWorkload, err
	}
	p.DirOut = "/chart"
	comp := compression.NewCompression()
	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
	if err != nil {
		return yr.TopologyBaseInfraWorkload, err
	}

	p.DirIn = "/chart"
	err = yr.ReadK8sWorkloadInMemFsDir(p, m)
	return yr.TopologyBaseInfraWorkload, err
}
