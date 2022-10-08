package v0

import (
	"io/ioutil"
	"os"

	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

func (l *Lib) CreateFile(p structs.Path, data []byte) error {
	// make path if it doesn't exist
	if _, err := os.Stat(p.FileOutPath()); os.IsNotExist(err) {
		_ = os.MkdirAll(p.DirIn, 0700) // Create your dir
	}
	return l.Log.ErrHandler(ioutil.WriteFile(p.FileOutPath(), data, 0644))
}

// OpenFile requires you to know that you need to close this
func (l *Lib) OpenFile(p structs.Path) (*os.File, error) {
	f, err := os.OpenFile(p.FileInPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	return f, l.Log.ErrHandler(err)
}

func (l *Lib) DeleteFile(p structs.Path) error {
	err := os.Remove(p.FileInPath())
	return l.Log.ErrHandler(err)
}
