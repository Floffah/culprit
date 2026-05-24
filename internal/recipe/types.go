package recipe

type Type string

const (
	TypeSoft    Type = "soft"
	TypeHard    Type = "hard"
	TypeCulprit Type = "culprit"
)

type Recipe struct {
	Name        string   `toml:"name" jsonschema:"description=The name of the recipe\\, must be unique across all recipes."`
	Description string   `toml:"description" jsonschema:"description=A brief description of the recipe."`
	Type        Type     `toml:"type" jsonschema:"description=The type of the recipe.,enum=soft,enum=hard,enum=culprit"`
	Targets     []Target `toml:"targets" jsonschema:"description=The targets of the recipe."`
}

type TargetKind string

const (
	TargetAbsolute TargetKind = "absolute"
)

type Target struct {
	Reason string     `toml:"reason" jsonschema:"description=The reason why this target should be cleaned up."`
	Kind   TargetKind `toml:"kind" jsonschema:"description=The kind of the target.,enum=absolute"`
	// Absolute path to the file to be cleaned up, supports glob patterns
	// Only allowed with "absolute" kind
	Path string `toml:"path" jsonschema:"description=Absolute path to the file to be cleaned up, supports glob patterns. Only allowed with \"absolute\" kind."`
}
