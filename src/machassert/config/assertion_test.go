package config

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSimpleAssertionParse(t *testing.T) {
	spec, err := ParseAssertionsSpecFile("testdata/assertions/simple.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec == nil {
		t.Fatal("Non-nil spec expected")
	}

	if spec.Name != "frontend" {
		t.Errorf("Got spec.Name=%q, wanted 'frontend'", spec.Name)
	}

	if spec.Assertions["binary"] == nil {
		t.Fatal("No assertion 'binary'")
	}
	a1 := spec.Assertions["binary"]
	if a1.Kind != FileExistsAssrt || a1.FilePath != "/bin/ls" {
		t.Errorf("Got %q, wanted {Kind:'exists', FilePath:'/bin/ls'}", spew.Sdump(a1))
	}

	if len(a1.Actions) != 1 {
		t.Fatalf("Got len(.Actions) == %d, wanted 1", len(a1.Actions))
	}
	if a1.Actions[0].Kind != "FAIL" {
		t.Errorf("Got kind=%q, wanted 'FAIL'", a1.Actions[0].Kind)
	}
}
