package main

import (
	"errors"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/starkriedesel/prompter"
)

func main() {
	// Basic command
	helloCmd := prompter.Command{
		Name:        "hello",
		Description: "print hello world",
		Executor: func(_ prompter.CmdArgs) error {
			fmt.Println("Hello World")
			return nil
		},
	}

	// Command with 1 sub command
	sayCmd := prompter.Command{
		Name:        "say",
		Description: "say some words",
	}
	fooSubCmd := prompter.Command{
		Name:        "foo",
		Description: "say foo",
		Executor: func(_ prompter.CmdArgs) error {
			fmt.Println("Prompter says \"foo\"")
			return nil
		},
	}
	sayCmd.AddSubCommands(fooSubCmd)

	// Greet Command
	nameArg := prompter.Argument{
		Name:              "--name",
		Description:       "your name",
		ArgumentCompleter: nameCompletor,
	}
	greetCmd := prompter.Command{
		Name:        "greet",
		Description: "say a greeting",
		Executor:    greetFunction,
	}
	greetCmd.AddArguments(nameArg)

	// Exit command
	exitCmd := prompter.ExitCommand("exit", "exit the application")

	// Create the prompter completer
	completer := prompter.NewCompleter()
	completer.RegisterCommands(helloCmd, sayCmd, greetCmd, exitCmd)

	// Start go-prompt
	p := prompt.New(completer.Execute, completer.Complete,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("pewpew"),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()
}

func greetFunction(args prompter.CmdArgs) error {
	if !args.Contains("--name") {
		return errors.New("must provide a name")
	}
	name, err := args.GetFirstValue("--name")
	if err != nil {
		return err
	}
	fmt.Printf("Hello to %s\n", name)
	return nil
}

func nameCompletor(_ string, _ []string) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "alice"},
		{Text: "bob"},
		{Text: "charles"},
	}
}
