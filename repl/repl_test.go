package repl

import (
	"strings"
	"testing"
)

func TestREPLClassAcrossLines(t *testing.T) {
	input := "" +
		"class A {\n" +
		"const new = function() {}\n" +
		"const cls.n = function() { return 7 }\n" +
		"}\n" +
		"var a = A.new()\n" +
		"a.n()\n"
	var output strings.Builder
	Start(strings.NewReader(input), &output)
	if !strings.Contains(output.String(), "7.000000") {
		t.Fatalf("repl output:\n%s", output.String())
	}
}

func TestREPLSemanticRejectsTypeChange(t *testing.T) {
	input := "var a = 1\na = \"no\"\n"
	var output strings.Builder
	Start(strings.NewReader(input), &output)
	if !strings.Contains(output.String(), "cannot be reassigned to another type") {
		t.Fatalf("repl output:\n%s", output.String())
	}
}
