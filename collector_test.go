package proxygen

import (
	"reflect"
	"testing"
)

func TestCollector_Compile(t *testing.T) {
	c := NewCollector()

	c.Load("./examples/interface.go")
	i, _ := c.FindInterface("SomeRepository")

	data, _ := i.Transform()

	got := data.Funcs

	want := []Func{
		{
			Name:         "Get",
			Receiver:     Var{"t", "*SomeRepositoryTracer"},
			Params:       []Var{{"id", "int"}},
			Return:       []Var{{"", "Some"}, {"", "error"}},
			OriginalType: "SomeRepository",
		},
		{
			Name:         "Save",
			Receiver:     Var{"t", "*SomeRepositoryTracer"},
			Params:       []Var{{"arg0", "Some"}},
			Return:       []Var{{"", "error"}},
			OriginalType: "SomeRepository",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nwant: %#v\n got: %#v", want, got)
	}
}
