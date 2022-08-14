package printer

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateFile(p, fn, folder string, data []byte) {
	// make path if it doesn't exist
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0700) // Create your dir
	}
	err := ioutil.WriteFile(p, data, 0644)
	if err != nil {
		log.Fatalf("error writing %s: %s", fn, err)
	}
}

// OpenFile requires you to know that you need to close this
func OpenFile(p string) (*os.File, error) {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalf("error writing %s: %s", f.Name(), err)
	}
	return f, err
}

func DeleteFile(p string) error {
	err := os.Remove(p)
	if err != nil {
		log.Fatalf("file %s deletion error %s", p, err.Error())
	}
	return err
}
