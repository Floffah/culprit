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
	Inputs      []Input  `toml:"inputs" jsonschema:"description=Inputs required by the recipe."`
	Targets     []Target `toml:"targets" jsonschema:"description=The targets of the recipe."`
}

type InputKind string

const (
	InputPath InputKind = "path"
)

type Input struct {
	ID          string    `toml:"id" jsonschema:"description=Stable identifier used to reference this input from targets."`
	Preset      string    `toml:"preset" jsonschema:"description=Named input preset used to populate kind\\, prompt\\, default\\, and global_id."`
	GlobalID    string    `toml:"global_id" jsonschema:"description=Shared identifier used to collect one value across recipes."`
	Kind        InputKind `toml:"kind" jsonschema:"description=The kind of input.,enum=path"`
	Prompt      string    `toml:"prompt" jsonschema:"description=Prompt shown when collecting this input."`
	Description string    `toml:"description" jsonschema:"description=Additional help text shown below the prompt."`
	Default     string    `toml:"default" jsonschema:"description=Default value used for this input."`
	Placeholder string    `toml:"placeholder" jsonschema:"description=Placeholder shown when the input is empty."`
}

type TargetKind string

const (
	TargetAbsolute TargetKind = "absolute"
	TargetCommand  TargetKind = "command"
)

type Target struct {
	Reason       string     `toml:"reason" jsonschema:"description=The reason why this target should be cleaned up."`
	Kind         TargetKind `toml:"kind" jsonschema:"description=The kind of the target.,enum=absolute,enum=command"`
	Path         string     `toml:"path" jsonschema:"description=Absolute path to the file to be cleaned up, supports glob patterns. Only allowed with \"absolute\" kind."`
	Command      string     `toml:"command" jsonschema:"description=The command to be executed to clean up the target. Only allowed with \"command\" kind."`
	RelatedPaths []string   `toml:"related_paths" jsonschema:"description=List of related paths to be used to calculate clearable size. Only allowed with \"command\" kind."`
}
