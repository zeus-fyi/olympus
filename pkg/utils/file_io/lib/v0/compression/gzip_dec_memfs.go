package compression

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func (c *Compression) UnGzipFromInMemFsOutToInMemFS(p *filepaths.Path, fs memfs.MemFS) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := fs.ReadFileInPath(p)
	if err != nil {
		log.Err(err).Msgf("Compression: UnGzipFromInMemFsOutToInMemFS, fs.ReadFileInPath(p) %s", p.FileInPath())
		return err
	}
	in := &bytes.Buffer{}
	_, err = in.Write(r)
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, herr := tr.Next()

		switch {

		// if no more files are found return
		case herr == io.EOF:
			return nil

		// return any other error
		case herr != nil:
			return herr

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// check the file type
		switch header.Typeflag {
		// if it's a dir do nothing
		case tar.TypeDir:
		// if it's a file create it
		case tar.TypeReg:
			p.FnIn = header.Name

			fo := p.FileDirOutFnInPath()
			p.DirIn = path.Dir(fo)

			out := &bytes.Buffer{}
			if _, cerr := io.Copy(out, tr); cerr != nil {
				log.Err(err).Msg("Compression: UnGzipFromInMemFsOutToInMemFS, io.Copy(out, tr)")
				return cerr
			}
			ferr := fs.MakeFileDirOutFnInPath(p, out.Bytes())
			if ferr != nil {
				log.Err(err).Msg("Compression: UnGzipFromInMemFsOutToInMemFS, fs.MakeFile(p, out.Bytes())")
				return ferr
			}
		}
	}
}
