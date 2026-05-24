package recipe

type Type string

const (
	TypeSoft      Type = "soft"
	TypeHard      Type = "hard"
	TypeMaculprit Type = "maculprit"
)

type Recipe struct {
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Type        Type     `toml:"type"`
	Targets     []Target `toml:"targets"`
}

type TargetKind string

const (
	TargetAbsolute TargetKind = "absolute"
)

type Target struct {
	Reason string     `toml:"reason"`
	Kind   TargetKind `toml:"kind"`
	// Absolute path to the file to be cleaned up, supports glob patterns
	// Only allowed with "absolute" kind
	Path string `toml:"path"`
}
