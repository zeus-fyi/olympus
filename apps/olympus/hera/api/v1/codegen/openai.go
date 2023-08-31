package ai_codegen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func CodeGenRoutes(e *echo.Group) *echo.Group {
	e.POST("/openai/codegen", CreateCodeGenAPIRequestHandler)
	e.POST("/ui/openai/codegen", CreateUICodeGenAPIRequestHandler)
	return e
}

func CreateUICodeGenAPIRequestHandler(c echo.Context) error {
	request := new(UICodeGenAPIRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CompleteUICodeGenRequest(c)
}
func (ai *UICodeGenAPIRequest) CompleteUICodeGenRequest(c echo.Context) error {
	ctx := context.Background()
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Ctx(ctx).Err(fmt.Errorf("failed to cast orgUser")).Msg("CompleteUICodeGenRequest")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp, err := hera_openai.HeraOpenAI.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: ai.Prompt,
					Name:    fmt.Sprintf("%d", ou.UserID),
				},
			},
		},
	)
	if err != nil {
		log.Ctx(ctx).Info().Interface("ou", ou).Interface("prompt", ai.Prompt).Interface("resp", resp).Err(err).Msg("CompleteUICodeGenRequest: CreateChatCompletion")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = hera_openai.HeraOpenAI.RecordUIChatRequestUsage(ctx, ou, resp)
	if err != nil {
		log.Ctx(ctx).Info().Interface("ou", ou).Interface("prompt", ai.Prompt).Interface("resp", resp).Err(err).Msg("CompleteUICodeGenRequest: RecordUIChatRequestUsage")
		//return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Choices[0].Message.Content)
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
	maxTokens := c.FormValue("maxTokens")
	tokens, err := strconv.Atoi(maxTokens)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusBadRequest, err)
	}
	params := hera_openai.OpenAIParams{
		Model:     model,
		MaxTokens: tokens,
		Prompt:    prompt,
	}
	cg, err := hera_openai.HeraOpenAI.MakeCodeGenRequest(ctx, ou, params)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, cg)
}

type UICodeGenAPIRequest struct {
	TokenEstimate int    `json:"tokenEstimate,omitempty"`
	Prompt        string `json:"prompt"`
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
