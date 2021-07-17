package proxygen

import (
	"fmt"
	"io"
	"text/template"
)

type (
	// Generator will generate code with input Proxy
	Generator struct {
		sg *structGenerator
		fg *funcGenerator
	}

	// funcGenerator will generate function code
	funcGenerator struct {
		tmpl *template.Template
	}

	// structGenerator will generate struct code
	structGenerator struct {
		tmpl *template.Template
	}
)

// NewGenerator will instatiate Generator that have FuncGenerator and StructGenerator
func NewGenerator() *Generator {
	return &Generator{
		sg: newStructGenerator(),
		fg: newFuncGenerator(),
	}
}

// GenerateAll with generat all proxies from Proxies into out (io.Writer)
func (g *Generator) GenerateAll(out io.Writer, data []Proxy) error {
	var err error

	if len(data) < 1 {
		return fmt.Errorf("proxies object is empty")
	}

	packageLine := fmt.Sprintf("package %s\n", data[0].PackageName)
	_, err = out.Write([]byte(packageLine))
	if err != nil {
		return err
	}

	for _, datum := range data {
		err = g.Generate(out, datum)
		if err != nil {
			return err
		}
	}

	return nil
}

// Generate will generate proxy struct and all its funs into out (io.Writer)
// this function doesnt generate package name at first line. If you need to
// also generate package you can use GenerateAll
func (g *Generator) Generate(out io.Writer, data Proxy) error {
	var err error

	err = g.sg.Generate(out, data)
	if err != nil {
		return err
	}

	for _, f := range data.Funcs {
		err = g.fg.generate(out, f)
		if err != nil {
			return err
		}
	}

	return err
}

// newFuncGenerator will instantiate FuncGenerator
func newFuncGenerator() *funcGenerator {
	f := &funcGenerator{}

	err := f.SetTemplate(defaultFuncTemplate)
	if err != nil {
		panic(fmt.Errorf("error parsing default template: %v", err))
	}

	return f
}

// SetTemplate is for
func (f *funcGenerator) SetTemplate(templateString string) error {
	tmpl, err := template.New("func").Parse(templateString)
	if err != nil {
		return err
	}

	f.tmpl = tmpl
	return nil
}

func (f *funcGenerator) generate(out io.Writer, input Func) error {
	data := input.FuncText()
	return f.tmpl.Execute(out, data)
}

const defaultFuncTemplate = `
func ({{ .ReceiverText }}) {{ .Name }} {{ .ParamsText }} {{ .ReturnText }} {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: {{.Receiver.Type}}.{{.Name}}: %v", dif)
	}(time.Now())
	return {{ .Receiver.Name }}.{{ .BaseType }}.{{ .Name }}{{ .ParamsNames }}
}
`

func newStructGenerator() *structGenerator {
	f := &structGenerator{}

	err := f.SetTemplate(defaultStructTemplate)
	if err != nil {
		panic(fmt.Errorf("error parsing default template: %v", err))
	}

	return f
}

func (f *structGenerator) SetTemplate(templateString string) error {
	tmpl, err := template.New("struct").Parse(templateString)
	if err != nil {
		return err
	}

	f.tmpl = tmpl
	return nil
}

//
func (f *structGenerator) Generate(out io.Writer, p Proxy) error {
	return f.tmpl.Execute(out, p)
}

const defaultStructTemplate = `
type {{ .Type }} struct {
	{{ .BaseType }}
}
`
