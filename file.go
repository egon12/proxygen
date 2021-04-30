package proxygen

import (
	"os"
	"os/exec"
	"strings"
)

func GenerateTracerFileName(filename string) string {
	return strings.Replace(filename, ".go", "_tracer.go", 1)
}

func FixImports(filename string) error {
	content, err := exec.Command("goimports", filename).CombinedOutput()
	if err != nil {
		return err
	}

	output, err := os.Create(filename)
	defer output.Close()

	if err != nil {
		return err
	}

	_, err = output.Write(content)
	if err != nil {
		return err
	}

	return nil
}
