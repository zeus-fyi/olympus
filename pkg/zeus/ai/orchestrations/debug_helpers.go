package ai_platform_service_orchestrations

import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func ChangeToAiDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

type CsvIteratorDebug struct {
	Cp *MbChildSubProcessParams
}

func (f *CsvIteratorDebug) Save() {
	dirMain := ChangeToAiDir()
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	ch := chronos.Chronos{}
	rn = "CsvIteratorDebug"
	fp := filepaths.Path{
		DirIn:  dirMain,
		DirOut: path.Join(dirMain, "tmp"),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset, ch.UnixTimeStampNow()),
	}
	err = fp.WriteToFileOutPath(b)
	if err != nil {
		panic(err)
	}
}

func OpenCsvIteratorDebug(fn string) CsvIteratorDebug {
	dirMain := ChangeToAiDir()
	f := CsvIteratorDebug{}
	if f.Cp == nil {
		f.Cp = &MbChildSubProcessParams{
			Wsr: artemis_orchestrations.WorkflowStageReference{
				RunCycle:    1,
				ChunkOffset: 0,
			},
		}
	}
	fp := filepaths.Path{
		DirIn:  path.Join(dirMain, "tmp"),
		DirOut: path.Join(dirMain, "tmp"),
		FnIn:   fn,
		FnOut:  fn,
	}
	b := fp.ReadFileInPath()
	err := json.Unmarshal(b, &f)
	if err != nil {
		panic(err)
	}
	return f
}

type AiAggregateAnalysisRetrievalTaskInputDebug struct {
	SourceTaskIds []int
	Cp            *MbChildSubProcessParams
}

type DebugJsonOutputs struct {
	Mb            *MbChildSubProcessParams
	Params        hera_openai.OpenAIParams
	JsonResponses []map[string]interface{}
}

func (f *DebugJsonOutputs) Save() {
	dirMain := ChangeToAiDir()
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	//ch := chronos.Chronos{}
	//ch.UnixTimeStampNow(),
	//rn := te
	fp := filepaths.Path{
		DirIn:  dirMain,
		DirOut: path.Join(dirMain, "tmp"),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Mb.Wsr.RunCycle, f.Mb.Wsr.ChunkOffset),
	}
	err = fp.WriteToFileOutPath(b)
	if err != nil {
		panic(err)
	}
}

var rn = "GenerateCycleReports"

func (f *AiAggregateAnalysisRetrievalTaskInputDebug) Save() {
	dirMain := ChangeToAiDir()
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		panic(err)
	}
	//ch := chronos.Chronos{}
	//ch.UnixTimeStampNow(),
	//rn := te
	fp := filepaths.Path{
		DirIn:  dirMain,
		DirOut: path.Join(dirMain, "tmp"),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset),
	}
	err = fp.WriteToFileOutPath(b)
	if err != nil {
		panic(err)
	}
}

func (f *AiAggregateAnalysisRetrievalTaskInputDebug) OpenFp() filepaths.Path {
	dirMain := ChangeToAiDir()
	//rn := "AiAggregateAnalysisRetrievalTaskInputDebug"
	if f.Cp == nil {
		f.Cp = &MbChildSubProcessParams{
			Wsr: artemis_orchestrations.WorkflowStageReference{
				RunCycle:    1,
				ChunkOffset: 0,
			},
		}
	}
	fp := filepaths.Path{
		DirIn:  path.Join(dirMain, "tmp"),
		DirOut: path.Join(dirMain, "tmp"),
		FnIn:   fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset),
	}
	return fp
}

type FanOutApiCallRequestTaskInputDebug struct {
	Rts []iris_models.RouteInfo
	Cp  *MbChildSubProcessParams
}

func (f *FanOutApiCallRequestTaskInputDebug) Open() {
	dirMain := ChangeToAiDir()
	rn := "FanOutApiCallRequestTaskInputDebug"
	fp := filepaths.Path{
		DirIn:  dirMain,
		DirOut: path.Join(dirMain, "tmp"),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset),
	}
	err := json.Unmarshal(fp.ReadFileInPath(), &f)
	if err != nil {
		panic(err)
	}
}

func (f *FanOutApiCallRequestTaskInputDebug) Save() {
	dirMain := ChangeToAiDir()
	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	//ch := chronos.Chronos{}
	//ch.UnixTimeStampNow(),
	rn := "FanOutApiCallRequestTaskInputDebug"
	fp := filepaths.Path{
		DirIn:  dirMain,
		DirOut: path.Join(dirMain, "tmp"),
		FnOut:  fmt.Sprintf("%s-cycle-%d-chunk-%d.json", rn, f.Cp.Wsr.RunCycle, f.Cp.Wsr.ChunkOffset),
	}
	err = fp.WriteToFileOutPath(b)
	if err != nil {
		panic(err)
	}
}
