package proxygen

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	input := Func{
		Name:     "Func1",
		Receiver: Var{"s", "*Struct1"},
		Params: []Var{
			{"args0", "Args0"},
			{"args1", "Args1"},
		},
		Return: []Var{
			{"", "Ret0"},
			{"", "error"},
		},
	}

	out := &strings.Builder{}
	g := NewFuncGenerator()
	g.Generate(out, input)
	got := out.String()

	want := `
func (s *Struct1) Func1 (args0 Args0,args1 Args1) ( Ret0, error) {
	defer func(start time.Time) {
		end := time.Now()
		dif := end.Sub(start)
		log.Printf("Duration: *Struct1.Func1: %v", dif)
	}(time.Now())
	return s.real.Func1(args0,args1)
}
`

	if got != want {
		t.Errorf("\nwant: %s\n got: %s", want, got)
		t.Errorf("\nwant: %v\n got: %v", []byte(want), []byte(got))
	}
}
