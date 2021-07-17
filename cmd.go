package proxygen

import (
	"fmt"
	"os"
	"strings"
)

// This is to used by cmd or main pacage
func Generate(filename, interfaceNames string) error {
	var err error

	c := NewCollector()
	err = c.Load(filename)
	if err != nil {
		return err
	}

	names := strings.Split(interfaceNames, ",")

	proxies := make([]Proxy, len(names))

	t := &transformer{}

	for i, name := range names {
		ci, err := c.FindInterface(name)
		if err != nil {
			return err
		}
		proxies[i], err = t.Transform(ci, "Tracer")
		if err != nil {
			return err
		}
	}

	outputFilename := GenerateTracerFileName(filename)
	out, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("error create %s: %v", outputFilename, err)
	}

	g := NewGenerator()
	err = g.GenerateAll(out, proxies)
	if err != nil {
		return err
	}

	_ = out.Close() // Close before fix imports
	err = FixImports(outputFilename)
	if err != nil {
		return err
	}

	return nil
}
