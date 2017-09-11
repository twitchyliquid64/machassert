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
	if a1.Kind != FileExistsAssrt || a1.FilePath != "/bin/ls" || a1.Order != 1 {
		t.Errorf("Got %q, wanted {Kind:'exists', FilePath:'/bin/ls', Order: 1}", spew.Sdump(a1))
	}

	if len(a1.Actions) != 1 {
		t.Fatalf("Got len(.Actions) == %d, wanted 1", len(a1.Actions))
	}
	if a1.Actions[0].Kind != "FAIL" {
		t.Errorf("Got kind=%q, wanted 'FAIL'", a1.Actions[0].Kind)
	}

	if spec.Assertions["thing"] == nil {
		t.Fatal("No assertion 'thing'")
	}
	a2 := spec.Assertions["thing"]
	if a2.Order != 1000 {
		t.Error("Expected default order '1000'")
	}
	t.Log(spew.Sdump(a2))
}

func TestSimpleAssertionErrors(t *testing.T) {
	_, err := ParseAssertionsSpecFile("testdata/assertions/badkind.hcl")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "unsupported assertion type/kind: welperino" {
		t.Errorf("Got %q, Want 'unsupported assertion type/kind: welperino'", err)
	}
}

func TestBadMatchFileAssertionErrors(t *testing.T) {
	_, err := ParseAssertionsSpecFile("testdata/assertions/badMatchFile.hcl")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "base_path/file_path must be specified for file_match assertions" {
		t.Errorf("Got %q, Want 'base_path/file_path must be specified for file_match assertions'", err)
	}
}

func TestBadHashMatchAssertionErrors(t *testing.T) {
	_, err := ParseAssertionsSpecFile("testdata/assertions/badHashMatch.hcl")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "hash/file_path must be specified for md5_match assertions" {
		t.Errorf("Got %q, Want 'hash/file_path must be specified for md5_match assertions'", err)
	}
}

func TestBadFileExistsAssertionErrors(t *testing.T) {
	_, err := ParseAssertionsSpecFile("testdata/assertions/badFileExists.hcl")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "file_path must be specified for exists and !exists assertions" {
		t.Errorf("Got %q, Want 'file_path must be specified for exists and !exists assertions'", err)
	}
}

func TestBadCopyActionErrors(t *testing.T) {
	_, err := ParseAssertionsSpecFile("testdata/actions/badCopy.hcl")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "source_path/destination_path must be specified for COPY actions" {
		t.Errorf("Got %q, Want 'source_path/destination_path must be specified for COPY actions'", err)
	}
}
