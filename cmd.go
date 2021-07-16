package proxygen

import (
	"fmt"
	"os"
	"strings"
)

func Generate(filename, interfaceNames string) error {
	var err error

	c := NewCollector()
	err = c.Load(filename)
	if err != nil {
		return err
	}

	names := strings.Split(interfaceNames, ",")

	proxies := make([]Proxy, len(names))

	//i := &InterfaceTransformer{}

	for i, name := range names {
		t, err := c.FindInterface(name)
		if err != nil {
			return err
		}
		proxies[i], err = t.Transform()
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
	err = g.Generate(out, proxies)
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
