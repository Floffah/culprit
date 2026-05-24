package theming

import (
	"os"

	"github.com/charmbracelet/x/term"
)

var initialised bool
var areWeTTY bool

func AreWeTTY() bool {
	if !initialised {
		areWeTTY = term.IsTerminal(os.Stdout.Fd())
		initialised = true
	}

	return areWeTTY
}
