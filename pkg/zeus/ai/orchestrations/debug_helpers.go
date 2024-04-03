package ai_platform_service_orchestrations

import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
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

type AiAggregateAnalysisRetrievalTaskInputDebug struct {
	SourceTaskIds []int
	Cp            *MbChildSubProcessParams
}

func (f *AiAggregateAnalysisRetrievalTaskInputDebug) Save() {
	dirMain := ChangeToAiDir()
	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	//ch := chronos.Chronos{}
	//ch.UnixTimeStampNow(),
	rn := "AiAggregateAnalysisRetrievalTaskInputDebug"
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

func (f *AiAggregateAnalysisRetrievalTaskInputDebug) Open() {
	dirMain := ChangeToAiDir()
	rn := "AiAggregateAnalysisRetrievalTaskInputDebug"
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
