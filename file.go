package proxygen

import (
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/imports"
)

func GenerateTracerFileName(filename string) string {
	return strings.Replace(filename, ".go", "_tracer.go", 1)
}

// FixImports will fix and formattion any missing error in the generated code
func FixImports(filename string) error {
	//return fixImportsOld(filename)
	return fixImportsNew(filename)
}

func fixImportsOld(filename string) error {
	content, err := exec.Command("goimports", filename).CombinedOutput()
	if err != nil {
		return err
	}

	output, err := os.Create(filename)
	defer func() {
		_ = output.Close()
	}()

	if err != nil {
		return err
	}

	_, err = output.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func fixImportsNew(filename string) error {
	content, err := imports.Process(filename, nil, nil)

	output, err := os.Create(filename)
	defer func() {
		_ = output.Close()
	}()

	if err != nil {
		return err
	}

	_, err = output.Write(content)
	if err != nil {
		return err
	}

	return nil
}
