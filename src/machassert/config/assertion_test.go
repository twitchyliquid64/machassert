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
  if a1.Kind != FileExistsAssrt || a1.FilePath != "/bin/frontend-d" {
    t.Errorf("Got %q, wanted {Kind:'exists', FilePath:'/bin/frontend-d'}", spew.Sdump(a1))
  }
}
