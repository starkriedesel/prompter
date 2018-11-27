# prompter
Helper for [go-prompt](https://github.com/c-bata/go-prompt/). Under heavy development.

Thanks to @parsiya for taking my pile of code and turning it into something useful. Check out his [borrowedtime](https://github.com/parsiya/borrowedtime) app that uses prompter. 

## Example
Main example is at [test/main.go](test/main.go).

## Usage
The package was not designed to work in a vacuum. It's meant to be used in conjunction with `go-prompt`.

### Creating Commands
Commands are created with `Command` struct.

```go
type Command struct {
	Name        string // Name of command (to be typed on cli)
	Description string // Description shown for command by the completer
	Executor    ExecutorFunc // Function that executes when command is run
	Hide        func() bool  // function to optionally hide a command from the completer
	SubCommands []Command    // List of completable subcommands
	Arguments   []Argument   // List of completable arguments
	Completer   *Completer
}
```

Commands must have a `Name` and should have a useful `Description`. Commands that implement the `Executor` field will be executable. Command can be optionally hidden by defining a `Hide` function.

```go
helloCmd := prompter.Command{
	Name: "hello",
	Description: "print hello world", 
	Executor: func(_ prompter.CmdArgs) error {
		fmt.Println("Hello World")
		return nil
  },
}
```

#### Exit

One useful pre-packaged command is "exit" for quitting the cli application (via `os.Exit(0)`):

```go
exitCmd := prompter.ExitCommand("exit", "exit the application")
```


### Starting the prompt loop
To connect commands to go-prompt, create a `prompter.Completer` object, call `RegisterCommands(...Command)`, and create the go-prompt object with the `completer.Execute` and `completer.Complete` properties,

```go
// Create the prompter completer
completer := prompter.NewCompleter()
completer.RegisterCommands(helloCmd, exitCmd)

// Start go-prompt
p := prompt.New(completer.Execute, completer.Complete,
  prompt.OptionPrefix(">>> "),
  prompt.OptionTitle("pewpew"),
  prompt.OptionPrefixTextColor(prompt.White),
)
p.Run()
```

### SubCommands
Subcommands are normal commands and added to main commands. Main commands usually do not do anything by themselves. One or more subcommands can be added to another command using the `AddSubCommands(...Command)` call.

```go
sayCmd := prompter.Command{
  Name: "say",
  Description: "say some words",
}

fooSubCmd := prompter.Command{
  Name: "foo",
  Description: "say foo",
  Executor: func(_ prompter.CmdArgs) error {
    fmt.Println("Prompter says \"foo\"")
    return nil
  },
}

sayCmd.AddSubCommands(fooSubCmd)
```

### Arguments
Commands can have zero or more arguments using the `Argument` struct.

```go
type Argument struct {
	Name            string // Name of argument (to be typed on cli)
	Description     string // Description shown for argument by the completer
	OptionCompleter CompleterFunc // Function used to complete argument fields
	Repeatable      bool // Boolean to allow argument to be provided multiple times with different values
}
```

The following code can be called on the command line like: `greet --name Bob`

```go
nameArg := prompter.Argument{
	Name: "--name",
	Description: "your name",
}

greetCmd := prompter.Command{
  Name: "greet",
  Description: "say a greeting",
  Executor: greetFunction,
}

greetCmd.AddArguments(nameArg)
```

### Executor Functions and CmdArgs
The `Executor` function of the `Command` struct has the following type

``` go
// CmdArgs contains the arguments passed to the command.
type CmdArgs map[string][]string

type ExecutorFunc func(CmdArgs) error
```

`CmdArgs` will contain the passed arguments and their value(s). This is where the subcommand is processed and action upon. From the above `greet` command example, the following defines the `greetFunction` which executes the command given the `--name` argument value.

```go
func greetFunction(args prompter.CmdArgs) error {
	if ! args.Contains("--name") {
		return errors.New("must provide a name")
	}
	name, err := args.GetFirstValue("--name")
	if err != nil {
		return err
	}
	fmt.Printf("Hello to %s\n", name)
	return nil
}
```

The `CmdArgs.GetFirstValue(name string)` function is a shortcut for `CmdArgs.GetValue(name string, n int)` which gets the nth occurance of the named argument when repeatable arguments are in use.

### Completers
Arguments can also have completers of type:

```go
type CompleterFunc func(string, []string) []prompt.Suggest
```

CompleterFunc is used to display suggestions for each argument. These suggestions can be hardcoded or dynamic (or a mix of both). The current argument is passed in `optName`. This allows multiple arguments to have the same completer. While there's only one option in this example, a `switch` is used to act as a blueprint for multiple arguments.

```go
endpointCompleter := func (optName string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	switch optName {
  case "-endpoint":
    // Hardcoded suggestion.
    sugs = append(sugs, prompt.Suggest{Text: "aws", Description: "Amazon endpoint"})
    
    // Do something to get some dynamic suggestions.
    // Assuming newEndpoints is a []string of endpoints.
    newEndpoints := GetEndpointsFromAPI()
    
    for _, ed := range newEndpoints {
        sugs = append(sugs,
      prompt.Suggest{Text: ed, Description: ed})
    }
  }
	return sugs
}

endpointArg := prompter.Argument{
	Name: "-endpoint",
	ArgumentCompleter: endpointCompleter,
}
```

## Licence
Opensourced under the Apache License v 2.0 license. See [LICENSE](LICENSE) for details.