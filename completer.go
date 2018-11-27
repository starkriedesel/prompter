package prompter

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

// Completer represents one go-prompt completer.
type Completer struct {
	Commands map[string]*Command
	Options  map[string]*Argument

	// cache for speed
	commandCache []prompt.Suggest
	optionsCache []prompt.Suggest
}

// NewCompleter creates and returns a new *Completer.
func NewCompleter() *Completer {
	return &Completer{
		Commands:     make(map[string]*Command),
		Options:      make(map[string]*Argument),
		commandCache: make([]prompt.Suggest, 0),
		optionsCache: make([]prompt.Suggest, 0),
	}
}

// RegisterCommands adds one or more top-level commands to the completer.
func (c *Completer) RegisterCommands(commands ...Command) error {
	for _, comm := range commands {
		if _, exists := c.Commands[comm.Name]; exists {
			return fmt.Errorf("command %s already exists", comm.Name)
		}
		if comm.Completer == nil {
			comm.Completer = NewCompleter()
			err := comm.Completer.RegisterCommands(comm.SubCommands...)
			if err != nil {
				return err
			}
			err = comm.Completer.RegisterOptions(comm.Arguments...)
			if err != nil {
				return err
			}
		}
		// If we don't do this, only pointer will be copied and last command
		// will overwrite everything.
		tempComm := comm
		c.Commands[comm.Name] = &tempComm
		c.commandCache = append(c.commandCache, prompt.Suggest{Text: comm.Name, Description: comm.Description})
	}
	return nil
}

// RegisterOptions adds one or more options to the completer.
func (c *Completer) RegisterOptions(opts ...Argument) error {
	for _, opt := range opts {
		if _, exists := c.Options[opt.Name]; exists {
			return fmt.Errorf("option %s already exists", opt.Name)
		}
		// Assign to a new variable and append, opt is not recreated in each iteration.
		tempOpt := opt
		c.Options[opt.Name] = &tempOpt
		c.optionsCache = append(c.optionsCache, prompt.Suggest{Text: opt.Name, Description: opt.Description})
	}
	return nil
}

// Execute parses the commands and arguments (TODO: Is this correct?).
func (c *Completer) Execute(in string) {
	args := splitCommandWords(in)
	found := c.ExecuteArgs(args)
	if !found {
		fmt.Printf("unknown or incomplete command: %s\n", args)
	}
}

func (c *Completer) ExecuteArgs(args []string) bool {
	commandName := ""
	if len(args) > 0 {
		commandName = args[0]
	}
	cmd, commandExists := c.Commands[commandName]

	// Execute the last command given.
	if commandExists {
		found := false
		if len(args) == 1 {
			found = cmd.Completer.ExecuteArgs([]string{})
		} else {
			found = cmd.Completer.ExecuteArgs(args[1:])
		}
		if found {
			return true
		}
		if cmd.Executor == nil {
			return false
		}
		var err error
		if len(args) == 1 {
			err = cmd.Executor(CmdArgs{})
		} else {
			err = cmd.Executor(cmd.Completer.collectArguments(args[1:]))
		}
		if err != nil {
			fmt.Printf("error encountered: %s\n", err.Error())
			// TODO: Change error color here.
			// Not sure how we can change color here unless add it as a field to
			// the completer.
			// color.Red("%s", err)
		}
		return true
	}
	return false
}

func (c *Completer) Complete(in prompt.Document) []prompt.Suggest {
	args := splitCommandWords(in.TextBeforeCursor())
	return c.CompleteArgs(args)
}

func (c *Completer) CompleteArgs(args []string) []prompt.Suggest {
	if len(args) == 1 { // suggestions on the root command
		return append(
			//prompt.FilterHasPrefix(c.commandCache, args[0], true),
			c.commandFilter(args, c.commandCache),
			c.optionFilter(args, c.optionsCache)...)
	}

	// Use sub-command completer if found
	subCmd := c.Commands[args[0]]
	if subCmd != nil {
		return subCmd.Completer.CompleteArgs(args[1:])
	}

	if len(args) >= 2 {
		optionName := args[len(args)-2]
		if option, exists := c.Options[optionName]; exists {
			if option.ArgumentCompleter == nil {
				return []prompt.Suggest{}
			}
			return c.optionValueFilter(args, option.ArgumentCompleter(optionName, args))
		}
	}

	return c.optionFilter(args, c.optionsCache)
}

func (c *Completer) commandFilter(args []string, suggestions []prompt.Suggest) []prompt.Suggest {
	suggestions = prompt.FilterHasPrefix(suggestions, args[0], true)
	var ret []prompt.Suggest
	for _, s := range suggestions {
		if command, exists := c.Commands[s.Text]; exists && (command.Hide == nil || !command.Hide()) {
			ret = append(ret, s)
		}
	}
	return ret
}

func (c *Completer) optionFilter(args []string, suggestions []prompt.Suggest) []prompt.Suggest {
	suggestions = prompt.FilterHasPrefix(suggestions, args[len(args)-1], true)
	var ret []prompt.Suggest
	for _, s := range suggestions {
		found := false
		for _, a := range args {
			if a != s.Text {
				continue
			}
			found = true
			if option, exists := c.Options[a]; exists {
				if option.Repeatable {
					ret = append(ret, s)
				}
			} else {
				ret = append(ret, s)
			}
			break
		}
		if !found {
			ret = append(ret, s)
		}
	}
	return ret
}

func (c *Completer) optionValueFilter(args []string, suggestions []prompt.Suggest) []prompt.Suggest {
	suggestions = prompt.FilterHasPrefix(suggestions, args[len(args)-1], true)
	for i, s := range suggestions {
		if strings.Contains(s.Text, " ") {
			suggestions[i].Text = "\"" + s.Text + "\""
		}
	}
	return suggestions
}

func (c *Completer) collectArguments(args []string) CmdArgs {
	var optName string
	arguments := CmdArgs{}
	for _, arg := range args {
		if optName != "" {
			arguments[optName] = append(arguments[optName], arg)
			optName = ""
		} else {
			if _, exists := c.Options[arg]; exists {
				optName = arg
			} else {
				arguments["_"] = append(arguments["_"], arg)
			}
		}
	}
	return arguments
}

func splitCommandWords(in string) []string {
	parts := strings.Split(in, " ")
	var output []string
	inQuotes := false
	for _, p := range parts {
		if inQuotes {
			// TODO: do better, this will still fail on weird input such as \\"
			if strings.HasSuffix(p, "\"") && !strings.HasSuffix(p, "\\\"") {
				p = strings.TrimSuffix(p, "\"")
				inQuotes = false
			}
			output[len(output)-1] += " " + p
		} else {
			if strings.HasPrefix(p, "\"") && !strings.HasPrefix(p, "\\\"") {
				p = strings.TrimPrefix(p, "\"")
				inQuotes = true
			}
			output = append(output, p)
		}
	}
	return output
}
