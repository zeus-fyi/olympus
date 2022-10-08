package v0

import "github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"

func (l *Lib) CleanUpPaths(paths ...structs.Path) error {
	for _, p := range paths {
		if err := l.Log.ErrHandler(l.DeleteFile(p)); err != nil {
			return err
		}
	}
	return nil
}

func (l *Lib) NewPkgPath(pkg, dir, fn string) structs.Path {
	return structs.Path{
		PackageName: pkg,
		Dir:         dir,
		Fn:          fn,
	}
}

func (l *Lib) NewPath(dir, fn string) structs.Path {
	return structs.Path{
		Dir: dir,
		Fn:  fn,
	}
}

func (l *Lib) NewFullPathDefinition(env, pkg, dir, fn string) structs.Path {
	return structs.Path{
		Env:         env,
		PackageName: pkg,
		Dir:         dir,
		Fn:          fn,
	}
}
