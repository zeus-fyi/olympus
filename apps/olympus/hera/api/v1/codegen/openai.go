package ai_codegen

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/hera/openai"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func CodeGenRoutes(e *echo.Group) *echo.Group {
	e.POST("/openai/codegen", CreateCodeGenAPIRequestHandler)
	return e
}

func CreateCodeGenAPIRequestHandler(c echo.Context) error {
	request := new(CodeGenAPIRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CompleteCodeGenRequest(c)
}

func (ai *CodeGenAPIRequest) CompleteCodeGenRequest(c echo.Context) error {
	file, err := c.FormFile("prompt")
	if err != nil {
		log.Err(err).Msg("CompleteCodeGenRequest: FormFile")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	src, err := file.Open()
	if err != nil {
		log.Err(err).Msg("CompleteCodeGenRequest: file.Open()")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	defer src.Close()
	in := bytes.Buffer{}
	if _, err = io.Copy(&in, src); err != nil {
		log.Err(err).Msg("CompleteCodeGenRequest: Copy")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	prompt, err := UnGzipTextFiles(&in)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	model := c.FormValue("model")
	cg, err := openai.HeraOpenAI.MakeCodeGenRequest(ctx, prompt, model, ou)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, cg)
}

type CodeGenAPIRequest struct {
	Prompt string `json:"prompt"`
}

func UnGzipTextFiles(in *bytes.Buffer) (string, error) {
	p := filepaths.Path{DirIn: "/tmp", DirOut: "/tmp", FnIn: "prompt.tar.gz"}
	m := memfs.NewMemFs()
	err := m.MakeFileIn(&p, in.Bytes())
	if err != nil {
		log.Err(err)
		return "", err
	}
	p.DirOut = "/prompt"
	comp := compression.NewCompression()
	err = comp.UnGzipFromInMemFsOutToInMemFS(&p, m)
	if err != nil {
		log.Err(err)
		return "", err
	}
	p.DirIn = "/prompt"
	return AppendFilesToString(m, p.DirOut)
}

func AppendFilesToString(fs memfs.MemFS, dir string) (string, error) {
	var buffer bytes.Buffer
	// Read the directory
	files, ferr := fs.ReadDir(dir)
	if ferr != nil {
		log.Err(ferr)
		return "", ferr
	}
	// Iterate through the files in the directory
	for _, file := range files {
		// Open the file
		f, err := fs.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Err(err)
			return "", err
		}
		// Read the file content
		content, err := io.ReadAll(f)
		if err != nil {
			log.Err(err)
			f.Close()
			return "", err
		}
		// Write the file content to the buffer
		_, err = buffer.Write(content)
		if err != nil {
			log.Err(err)
			f.Close()
			return "", err
		}
		f.Close()
	}
	return buffer.String(), nil
}
