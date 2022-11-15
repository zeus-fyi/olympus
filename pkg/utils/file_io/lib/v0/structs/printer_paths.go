package structs

import (
	"path"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Path struct {
	PackageName string
	DirIn       string
	DirOut      string
	FnIn        string
	FnOut       string
	Env         string
	FilterFiles string_utils.FilterOpts
}

type Paths struct {
	Slice []Path
}

func (ps *Paths) AddPathToSlice(p Path) {
	ps.Slice = append(ps.Slice, p)
}

func (p *Path) FileDirOutFnInPath() string {
	return path.Join(p.DirOut, p.FnIn)
}

func (p *Path) FileInPath() string {
	return path.Join(p.DirIn, p.FnIn)
}

func (p *Path) FileOutPath() string {
	return path.Join(p.DirOut, p.FnOut)
}

func (p *Path) LeftExtendDirInPath(dirExtend string) string {
	p.DirIn = path.Join(dirExtend, p.DirIn)
	return p.DirIn
}

func (p *Path) RightExtendDirInPath(dirExtend string) string {
	p.DirIn = path.Join(p.DirIn, dirExtend)
	return p.DirIn
}

func (p *Path) LeftExtendDirOutPath(dirExtend string) string {
	p.DirOut = path.Join(dirExtend, p.DirOut)
	return p.DirOut
}

func (p *Path) RightExtendDirOutPath(dirExtend string) string {
	p.DirOut = path.Join(p.DirOut, dirExtend)
	return p.DirOut
}

func (p *Path) Local() string {
	return "artifacts/local/"
}

func (p *Path) Dev() string {
	return "artifacts/dev/"
}

func (p *Path) Staging() string {
	return "artifacts/staging/"
}

func (p *Path) Production() string {
	return "artifacts/production/"
}

func (p *Path) AddGoFn(fn string) {
	p.FnIn = fn + ".go"
}
