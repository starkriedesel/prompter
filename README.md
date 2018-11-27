# prompter
Helper for [go-prompt](https://github.com/c-bata/go-prompt/). Under heavy development.

## Example
Main example is at [test/main.go](test/main.go).

## Usage
The package was not designed to work in a vacuum. It's meant to be used in conjunction with `go-prompt`.

### Creating Commands
Commands are created with `Command` or `SubCommand`. They are both of type `Cmd`.

Each command or subcommands has at least a `Name` and a `Description`. Both are displayed to user in the resulting prompt.

``` go
Cmd1 := prompter.Command("deploy command", "deploy the application")
```

### SubCommands
Subcommands are normal commands and added to main commands. Main commands usually do not do anything by themselves. If a main command needs to perform some action, create it as a `SubCommand`.

``` go
CloudSubCmd := prompter.SubCommand(
    "cloud",
    "deploys the cloud",
    cloudExecutor,
)
```

### Options
Each subcommand can have multiple options. These are added with `AddOption`.

``` go
CloudSubCmd.AddOption("-endpoint", "endpoint to deploy the cloud", false, endpointCompleter)
```

`false` as third argument means the option is not repeatable.

### Completers
`AddOption` can have an optional completer of type `CompleterFunc`:

``` go
type CompleterFunc func(string, []string) []prompt.Suggest
```

CompleterFunc is used to display suggestions for each argument. These suggestions can be hardcoded or dynamic (or a mix of both). The current argument is passed in `optName`. This allows multiple arguments to have the same completer. While there's only one option in this example, a `switch` is used to act as a blueprint for multiple arguments.

``` go
func endpointCompleter(optName string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	switch optName {
    case "-endpoint":
        // Hardcoded suggestion.
        sugs = append(sugs, prompt.Suggest{Text: "aws", "Amazon endpoint"})
        
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
```

### Executor Functions
`cloudExecutor` is of type `ExecutorFunc` or:

``` go
ExecutorFunc func(CmdArgs) error
```

Where `CmdArgs` is a `map[string][]string` that contains the passed options/arguments and their value(s). This is where the subcommand is processed and action upon.

``` go
func cloudExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside cloudExecutor")
	fmt.Printf("cloud deploy args: %v\n", args)

	if args.Contains("-endpoint") {
        // Do the deployment.
        err := cloudDeploy.DeployTheCloud(endpoint)
        if err != nil {
            return err
        }
        // Or "return cloudDeploy.DeployTheCloud(endpoint)"
	}
	return nil
}
```
