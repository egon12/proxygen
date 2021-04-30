package proxygen

import "testing"

func TestGenerateTracerFileName(t *testing.T) {
	input := "./examples/interface.go"

	want := "./examples/interface_tracer.go"

	got := GenerateTracerFileName(input)

	if got != want {
		t.Errorf("\nwant: %s\n got:%s", want, got)
	}
}
