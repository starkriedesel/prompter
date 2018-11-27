package prompter

import (
	"fmt"
	"os"
)

type ExecutorFunc func(CmdArgs) error

type Cmd struct {
	Name        string
	Description string
	Completer   *Completer
	Executor    ExecutorFunc
	Hide        func() bool
	Hidden      bool
	SubCommands []Cmd
	Options     []Option
}

// Command returns a new command. For subcommands please use NewSubCommand.
// SubCommands are optional and can be added with AddSubCommands later.
func Command(name, desc string, subs ...Cmd) Cmd {
	return Cmd{
		Name:        name,
		Description: desc,
		SubCommands: subs,
	}
}

// AddSubCommands adds one or more subcommands to the higher-level command.
func (c *Cmd) AddSubCommands(subs ...Cmd) error {
	if len(subs) == 0 {
		return fmt.Errorf("no subcommands provided")
	}
	lenBefore := len(c.SubCommands)
	c.SubCommands = append(c.SubCommands, subs...)

	// TODO: Remove this sanity check later.
	lenAfter := len(c.SubCommands)
	if lenAfter-lenBefore != len(subs) {
		// Not everything was added.
		return fmt.Errorf("added %d subcommands instead of %d",
			lenAfter-lenBefore, len(subs))
	}
	return nil
}

// SubCommand returns a new subcommand. Add it to a higher-level command
// with AddSubCommand.
func SubCommand(name, desc string, exec ExecutorFunc, options ...Option) Cmd {
	cmd := Command(name, desc)
	cmd.Executor = exec
	cmd.Options = append(cmd.Options, options...)
	return cmd
}

// AddOption creates an option and adds it to the command.
// TODO: Check for duplicate items or just overwrite? I think overwrite.
func (c *Cmd) AddOption(name, desc string, repeatable bool, comp CompleterFunc) {
	opt := newOption(name, desc, repeatable, comp)
	c.Options = append(c.Options, opt)
}

// ExitCommand returns a command that exits the application.
func ExitCommand(name, desc string) Cmd {
	return SubCommand(name, desc, exitExecutor)
}

// exitExecutor exits the application.
func exitExecutor(args CmdArgs) error {
	os.Exit(1)
	// Does this really matter?
	return nil
}
