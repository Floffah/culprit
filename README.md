# Culprit

Culprit is a command line tool that helps you delete common and known files that fill up your Mac's storage. It is simple, extendable, open, and secure.

Features:
- **Culprit will never actually delete files**, it generates a bash script with `rm` calls that you can review and run yourself
- **Highly extendable**, define rules with [recipes](#recipes)
- **Comes with a set of built-in recipes** for common apps and file types
- **Recipes are smart**, they can use advanced conditions like only deleting certain files if the app is not installed, or only deleting files older than a certain date
- **Ship recipes with your app** to help users clean up files related to your app
- **Open source**, contribute your own recipes and improvements to the project!
- **Secure**, all file paths are sanitized to prevent command injection (and it doesn't delete files itself anyway)

## Installation

As culprit is currently in early development, there are no pre-built binaries available.

You can build it yourself:

```bash
git clone https://github.com/floffah/culprit.git
cd culprit
go build -o culprit cmds/culprit/culprit.go
mv culprit /usr/local/bin/culprit
```

OR install it with go

```bash
go install github.com/floffah/culprit/cmds/culprit@latest
```

## Usage

Everything is done with `culprit clean [type]`.

Current types are:
- `soft`: deletes common files that are generally safe to delete, such as caches and logs
- `hard`: deletes more files that may be less safe to delete, such as app support files that are no longer needed and old downloads

Options:
- `--force`: overwrites the generated script rather than failing if it exists already

## Recipes

Recipes are define in TOML files and live in the recipes directory (default `~/.culprit/recipes`). You can add your own recipes here, or contribute them to the project! The recipe dir is configurable by editing `~/.culprit/culprit.toml` or with the `CULPRIT_RECIPES` environment variable.

Built in recipes are not stored in the recipe dir but are embedded in the binary, but they can be overridden by adding a recipe with the same name to the recipe dir.

Recipes have a few required fields:
- `name`: the name of the recipe, used for logging and debugging
- `type`: the type of the recipe, either `soft` or `hard`
- `description`: a description of the recipe, used for logging and debugging

Recipes then have two collection: targets and conditions.

Targets define the files that will be deleted if the recipe is run. You can view [the types file](./internal/recipe/types.go) for the full list of target types, eventually I will document them here.

Conditions define the conditions that must be met for the recipe to be run. You can view [the types file](./internal/recipe/types.go) for the full list of condition types, eventually I will document them here.
