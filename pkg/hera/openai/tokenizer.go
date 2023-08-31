package hera_openai

import (
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetTokenApproximate(prompt string) int {
	ForceDirToPythonDir()
	cmd := exec.Command("python", "tokenizer.py", prompt)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Err(err)
		return 0
	}
	// Convert the output to a string and split it by newline
	outputStr := string(output)
	outputLines := strings.Split(outputStr, "\n")

	// Convert the first line of the output (which should be the result) to an integer
	result, err := strconv.Atoi(outputLines[0])
	if err != nil {
		log.Err(err)
		return 0
	}
	return result
}

func ForceDirToPythonDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
