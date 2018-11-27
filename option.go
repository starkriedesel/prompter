package prompter

import (
	prompt "github.com/c-bata/go-prompt"
)

// Option struct and utilities.

// CompleterFunc returns completion options for commands.
type CompleterFunc func(string, []string) []prompt.Suggest

// Option represents arguments for a subcommand.
type Option struct {
	Name            string
	Description     string
	OptionCompleter CompleterFunc
	Repeatable      bool
}

// newOption creates a new command option.
func newOption(name, desc string, repeatable bool, comp CompleterFunc) Option {
	return Option{
		Name:            name,
		Description:     desc,
		Repeatable:      repeatable,
		OptionCompleter: comp,
	}
}
