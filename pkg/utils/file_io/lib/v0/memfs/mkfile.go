package memfs

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (m *MemFS) MakeFileDirOutFnInPath(p *filepaths.Path, content []byte) error {
	merr := m.MkPathDirAll(p)
	if merr != nil {
		return merr
	}
	if err := m.WriteFile(p.FileDirOutFnInPath(), content, 0644); err != nil {
		return err
	}
	return nil
}

func (m *MemFS) MakeFileIn(p *filepaths.Path, content []byte) error {
	merr := m.MkPathDirAll(p)
	if merr != nil {
		log.Err(merr).Msgf("MemFS, MakeFile fileIn path %s, fileOut path %s", p.FileInPath(), p.FileOutPath())
		return merr
	}
	fnPath := p.FileInPath()
	if err := m.WriteFile(fnPath, content, 0644); err != nil {
		log.Err(err).Msgf("MemFS, WriteFile, fileOut path %s", p.FileInPath())
		return err
	}
	return nil
}

func (m *MemFS) MakeFileOut(p *filepaths.Path, content []byte) error {
	merr := m.MkPathDirAll(p)
	if merr != nil {
		log.Err(merr).Msgf("MemFS, MakeFile fileIn path %s, fileOut path %s", p.FileInPath(), p.FileOutPath())
		return merr
	}
	if err := m.WriteFile(p.FileOutPath(), content, 0644); err != nil {
		log.Err(err).Msgf("MemFS, WriteFile, fileOut path %s", p.FileOutPath())
		return err
	}
	return nil
}
