package paths

import "path"

type Path struct {
	FileName string
	Dir      string
	Path     string
	Env      string
}

func JoinDirToFileName(dir, fn string) string {
	p := path.Join(dir, fn)
	return p
}
