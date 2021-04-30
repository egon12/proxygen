package proxygen

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

type (
	Func struct {
		Name     string
		Receiver Var
		Params   MultiVar
		Return   MultiVar
	}

	Var struct {
		Name string
		Type string
	}

	MultiVar []Var

	Proxy struct {
		PackageName  string
		Receiver     Var
		OriginalType string
		Funcs        []Func
	}

	FuncText struct {
		Func
		ReceiverText string
		ParamsText   string
		ReturnText   string
		ParamsNames  string
	}

	Generator struct {
		sg *StructGenerator
		fg *FuncGenerator
	}

	FuncGenerator struct {
		tmpl *template.Template
	}

	StructGenerator struct {
		tmpl *template.Template
	}
)

func NewGenerator() *Generator {
	return &Generator{
		sg: NewStructGenerator(),
		fg: NewFuncGenerator(),
	}
}

func (g *Generator) Generate(out io.Writer, data []Proxy) error {
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
		err = g.generateSingle(out, datum)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateSingle(out io.Writer, data Proxy) error {
	var err error

	err = g.sg.Generate(out, data)
	if err != nil {
		return err
	}

	for _, f := range data.Funcs {
		err = g.fg.Generate(out, f)
		if err != nil {
			return err
		}
	}

	return err
}

func NewFuncGenerator() *FuncGenerator {
	f := &FuncGenerator{}

	err := f.SetTemplate(defaultFuncTemplate)
	if err != nil {
		panic(fmt.Errorf("error parsing default template: %v", err))
	}

	return f
}

func (f *FuncGenerator) SetTemplate(templateString string) error {
	tmpl, err := template.New("func").Parse(templateString)
	if err != nil {
		return err
	}

	f.tmpl = tmpl
	return nil
}

func (f *FuncGenerator) Generate(out io.Writer, input Func) error {

	data := FuncText{
		Func:         input,
		ReceiverText: input.Receiver.Text(),
		ParamsText:   input.Params.Text(),
		ReturnText:   input.Return.Text(),
		ParamsNames:  input.Params.Names(),
	}

	return f.tmpl.Execute(out, data)
}

const defaultFuncTemplate = `
func ({{ .ReceiverText }}) {{ .Name }} {{ .ParamsText }} {{ .ReturnText }} {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: {{.Receiver.Type}}.{{.Name}}: %v", dif)
	}(time.Now())
	return {{ .Receiver.Name }}.real.{{ .Name }}{{ .ParamsNames }}
}
`

func NewStructGenerator() *StructGenerator {
	f := &StructGenerator{}

	err := f.SetTemplate(defaultStructTemplate)
	if err != nil {
		panic(fmt.Errorf("error parsing default template: %v", err))
	}

	return f
}

func (f *StructGenerator) SetTemplate(templateString string) error {
	tmpl, err := template.New("struct").Parse(templateString)
	if err != nil {
		return err
	}

	f.tmpl = tmpl
	return nil
}

func (f *StructGenerator) Generate(out io.Writer, p Proxy) error {
	return f.tmpl.Execute(out, p)
}

const defaultStructTemplate = `
type {{ .Receiver.Type }} struct {
	real {{ .OriginalType }}
}
`

func (v Var) Text() string {
	return v.Name + " " + v.Type
}

func (m MultiVar) Text() string {
	vars := make([]string, len(m))
	for i, v := range m {
		vars[i] = v.Text()
	}

	return "(" + strings.Join(vars, ",") + ")"
}

func (m MultiVar) Names() string {
	vars := make([]string, len(m))
	for i, v := range m {
		vars[i] = v.Name
	}

	return "(" + strings.Join(vars, ",") + ")"
}