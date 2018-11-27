package prompter

import (
	"os"
)

type ExecutorFunc func(CmdArgs) error

type Command struct {
	Name        string // Name of command (to be typed on cli)
	Description string // Description shown for command by the completer
	Completer   *Completer
	Executor    ExecutorFunc // Function that executes when command is run
	Hide        func() bool  // function to optionally hide a command from the completer
	SubCommands []Command    // List of completable subcommands
	Arguments   []Argument   // List of completable arguments
}

// AddSubCommands adds one or more subcommands to the higher-level command.
func (c *Command) AddSubCommands(subs ...Command) {
	c.SubCommands = append(c.SubCommands, subs...)
}

// AddOption creates an option and adds it to the command.
func (c *Command) AddArguments(opts ...Argument) {
	c.Arguments = append(c.Arguments, opts...)
}

// ExitCommand returns a command that exits the application.
func ExitCommand(name, desc string) Command {
	return Command{
		Name:        name,
		Description: desc,
		Executor:    exitExecutor,
	}
}

// exitExecutor exits the application.
func exitExecutor(_ CmdArgs) error {
	os.Exit(0)
	// Does this really matter?
	return nil
}
