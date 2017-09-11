package config

// Assertion kinds
const (
	FileExistsAssrt    string = "exists"
	FileNotExistsAssrt string = "!exists"
	HashMatchAssrt     string = "md5_match"
	HashFileAssrt      string = "file_match"
)

// Action kinds
const (
	ActionFail     string = "FAIL"
	ActionCopyFile string = "COPY"
)

// AssertionSpec describes the high-level schema for a file containing assertions.
type AssertionSpec struct {
	Name       string
	Assertions map[string]*Assertion `hcl:"assert"`
}

// Assertion describes the schema for a assertion.
type Assertion struct {
	Kind  string
	Order int

	// FileExistsAssrt & FileNotExistsAssrt & HashMatchAssrt
	FilePath string `hcl:"file_path"`
	// HashMatchAssrt
	Hash string //hex-encoded hash bytes

	// HashFileAssrt
	BasePath string `hcl:"base_path"`

	Actions []*Action `hcl:"or"`
}

// Action represents the schema for an action taken on assertion failure.
type Action struct {
	Kind            string `hcl:"action"`
	SourcePath      string `hcl:"source_path"`
	DestinationPath string `hcl:"destination_path"`
}
