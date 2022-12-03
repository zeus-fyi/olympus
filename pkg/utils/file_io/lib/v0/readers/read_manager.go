package readers

import (
	"errors"
	"io/ioutil"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (l *ReaderLib) ReadFilePathPtr(p *filepaths.Path) []byte {
	if p == nil {
		panic(errors.New("no path provided"))
	}

	byteArray, err := ioutil.ReadFile(p.FileDirOutFnInPath())
	if err != nil {
		panic(err)
	}
	return byteArray
}
func (l *ReaderLib) ReadFile(p filepaths.Path) []byte {
	byteArray, err := ioutil.ReadFile(p.FileDirOutFnInPath())
	if err != nil {
		panic(err)
	}
	return byteArray
}

func (l *ReaderLib) ReadJsonObject(p filepaths.Path) []byte {
	jsonByteArray, err := p.ReadFileInPath()
	if err != nil {
		panic(err)
	}

	return jsonByteArray
}

func (l *ReaderLib) ReadFilePathOutJsonObject(p filepaths.Path) []byte {
	jsonByteArray, err := p.ReadFileOutPath()
	if err != nil {
		panic(err)
	}

	return jsonByteArray
}
