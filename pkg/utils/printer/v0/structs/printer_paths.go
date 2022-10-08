package structs

import (
	"path"
)

type Path struct {
	PackageName string
	Dir         string
	Fn          string
	Env         string
}

type Paths []Path

func (p *Path) FilePath() string {
	return path.Join(p.Dir, p.Fn)
}

func (p Path) Local() string {
	return "artifacts/local/"
}

func (p Path) Dev() string {
	return "artifacts/dev/"
}

func (p Path) Staging() string {
	return "artifacts/staging/"
}

func (p Path) Production() string {
	return "artifacts/production/"
}
