package prompter

import (
	"github.com/c-bata/go-prompt"
)

// Argument struct and utilities.

// CompleterFunc returns completion options for commands.
type CompleterFunc func(string, []string) []prompt.Suggest

// Argument represents arguments for a subcommand.
type Argument struct {
	Name              string
	Description       string
	ArgumentCompleter CompleterFunc
	Repeatable        bool
}
