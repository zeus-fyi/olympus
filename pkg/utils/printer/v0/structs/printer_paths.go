package structs

import (
	"path"
)

type Path struct {
	PackageName string
	DirIn       string
	DirOut      string
	Fn          string
	Env         string
}

type Paths struct {
	Slice []Path
}

func (ps *Paths) AddPathToSlice(p Path) {
	ps.Slice = append(ps.Slice, p)
}

func (p *Path) FileOutPath() string {
	return path.Join(p.DirOut, p.Fn)
}

func (p *Path) FileInPath() string {
	return path.Join(p.DirIn, p.Fn)
}

func (p *Path) LeftExtendDirInPath(dirExtend string) string {
	p.DirIn = path.Join(dirExtend, p.DirIn)
	return p.DirIn
}

func (p *Path) LeftExtendDirOutPath(dirExtend string) string {
	p.DirOut = path.Join(dirExtend, p.DirOut)
	return p.DirOut
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
