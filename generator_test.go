package proxygen

import (
	"strings"
	"testing"
)

func TestFuncGenerate(t *testing.T) {
	input := Func{
		Name:     "Func1",
		Receiver: Var{"s", "*Struct1Tracer"},
		Params: []Var{
			{"args0", "Args0"},
			{"args1", "Args1"},
		},
		Return: []Var{
			{"", "Ret0"},
			{"", "error"},
		},
		BaseType: "Struct1",
	}

	out := &strings.Builder{}
	g := newFuncGenerator()
	_ = g.generate(out, input)
	got := out.String()

	want := `
func (s *Struct1Tracer) Func1 (args0 Args0,args1 Args1) ( Ret0, error) {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *Struct1Tracer.Func1: %v", dif)
	}(time.Now())
	return s.Struct1.Func1(args0,args1)
}
`

	if got != want {
		t.Errorf("\nwant: %s\n got: %s", want, got)
		t.Errorf("\nwant: %v\n got: %v", []byte(want), []byte(got))
	}
}

func TestStructGenerator(t *testing.T) {
	input := Proxy{
		PackageName: "newpkg",
		Type:        "MyTracer",
		Funcs:       []Func{},
		BaseType:    "My",
	}

	out := &strings.Builder{}
	g := newStructGenerator()
	_ = g.Generate(out, input)
	got := out.String()

	want := `
type MyTracer struct {
	My
}
`

	if got != want {
		t.Errorf("\nwant: %s\n got: %s", want, got)
		t.Errorf("\nwant: %v\n got: %v", []byte(want), []byte(got))
	}

}
