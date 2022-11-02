package compression

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

// UnGzip takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func UnGzip(p structs.Path) error {
	r, err := os.Open(p.Fn)
	if err != nil {
		return err
	}
	defer r.Close()
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if merr := os.Mkdir(header.Name, 0755); merr != nil {
				return merr
			}
		// if it's a file create it
		case tar.TypeReg:
			p.Fn = header.Name

			fo := p.FileOutPath()
			b := path.Dir(fo)
			if _, zerr := os.Stat(b); os.IsNotExist(zerr) {
				_ = os.MkdirAll(b, 0700) // Create your dir
			}
			outFile, perr := os.Create(fo)
			if perr != nil {
				return perr
			}
			if _, cerr := io.Copy(outFile, tr); cerr != nil {
				return cerr
			}
			outFile.Close()
		}
	}
}
