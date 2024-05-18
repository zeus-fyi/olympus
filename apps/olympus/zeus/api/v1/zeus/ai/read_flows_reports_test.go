package zeus_v1_ai

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/memfs"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func (t *FlowsWorkerTestSuite) TestFlowExport() {
	tmpOu := t.Ou
	tmpOu.OrgID = 1685378241971196000

	gr, err := ExportRunCsvRequest2(ctx, tmpOu, 1715985973861270000)
	t.Require().Nil(err)

	inf := memfs.NewMemFs()
	// Create a new zip archive.
	for _, ue := range gr.Entities {
		for _, v := range ue.MdSlice {
			if v.TextData == nil {
				continue
			}
			p := filepaths.Path{
				FnIn:  ue.Nickname,
				DirIn: "/exports",
			}

			if v.TextData != nil {
				err = inf.MakeFileIn(&p, []byte(*v.TextData))
				t.Require().Nil(err)
			}
		}
	}

	p := filepaths.Path{
		FnIn:  gr.Name,
		DirIn: "/exports",
	}
	b, err := compression.GzipDirectoryToMemoryFS(p, inf)
	t.Require().Nil(err)

	p.FnOut = "tmp.tar.gz"
	p.DirOut = "/Users/alex/go/Olympus/olympus/apps/olympus/zeus/api/v1/zeus/ai"
	p.WriteToFileOutPath(b)
}

func TestFlowExportTestSuite(t *testing.T) {
	suite.Run(t, new(FlowsWorkerTestSuite))
}
