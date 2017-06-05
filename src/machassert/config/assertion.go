package config

// Assertion kinds
const (
	FileExistsAssrt    string = "exists"
	FileNotExistsAssrt string = "!exists"
	HashMatchAssrt     string = "md5_match"
)

// Action kinds
const (
	ActionFail string = "FAIL"
)

// AssertionSpec describes the high-level schema for a file containing assertions.
type AssertionSpec struct {
	Name       string
	Assertions map[string]*Assertion `hcl:"assert"`
}

// Assertion describes the schema for a assertion.
type Assertion struct {
	Kind string

	// FileExistsAssrt & FileNotExistsAssrt
	FilePath string `hcl:"file_path"`
	// HashMatchAssrt
	Hash string //hex-encoded hash bytes

	Actions []*Action `hcl:"or"`
}

type Action struct {
	Kind string `hcl:"action"`
}