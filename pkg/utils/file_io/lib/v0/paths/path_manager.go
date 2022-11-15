package paths

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (l *PathLib) CleanUpPaths(paths ...structs.Path) error {
	for _, p := range paths {
		if err := l.Log.ErrHandler(l.DeleteFile(p)); err != nil {
			return err
		}
	}
	return nil
}

func (l *PathLib) NewPkgPath(pkg, dir, fn string) structs.Path {
	return structs.Path{
		PackageName: pkg,
		DirIn:       dir,
		FnIn:        fn,
	}
}

func (l *PathLib) NewPath(dir, fn string) structs.Path {
	return structs.Path{
		DirIn: dir,
		FnIn:  fn,
	}
}

func (l *PathLib) NewPkgPathInOut(pkgName, dirIn, dirOut, fn string) structs.Path {
	return structs.Path{
		PackageName: pkgName,
		DirIn:       dirIn,
		DirOut:      dirOut,
		FnIn:        fn,
	}
}
func (l *PathLib) NewFullPathDefinition(env, pkg, dirIn, dirOut, fn string) structs.Path {
	return structs.Path{
		Env:         env,
		PackageName: pkg,
		DirIn:       dirIn,
		DirOut:      dirOut,
		FnIn:        fn,
	}
}
