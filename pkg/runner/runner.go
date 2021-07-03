package runner

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	template "github.com/zoomoid/papeer/pkg/template"
	types "github.com/zoomoid/papeer/pkg/types"
)

// Run creates a temporary directory and runs the templater to write to the file before calling latexmk on the file
func Run(options *types.WrapperOptions) (string, error) {
	dir, err := os.MkdirTemp("", "runner-*")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(dir)

	pwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	// redefine pdf to embed filename to qualified url
	options.Filename = filepath.Join(pwd, options.Filename)

	filename := uuid.New().String()

	contents, err := template.Template(options)

	if err != nil {
		return "", err
	}

	fname := filepath.Join(dir, filename+".tex")

	err = os.WriteFile(fname, []byte(contents), 0644)

	if err != nil {
		return "", err
	}

	cmd := exec.Command("latexmk", "-interaction=nonstopmode", "-pdf", "-shell-escape", "-outdir"+dir, "-synctex=0", "-file-line-error", fname)
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	tempOutputFilename := filepath.Join(dir, filename+".pdf")
	outputFilename := filepath.Join(pwd, filename+".pdf")

	err = os.Rename(tempOutputFilename, outputFilename)

	if err != nil {
		return "", err
	}

	return outputFilename, nil
}
