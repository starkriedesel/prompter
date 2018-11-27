package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/starkriedesel/prompter"
)

func main() {

	comp := prompter.NewCompleter()

	myCmd1 := prompter.Command("mycmd1", "mycmd1 description")

	subCmd1 := prompter.SubCommand("subcmd1", "subcmd1 description", subcmd1Executor)
	subCmd1.AddOption("-arg1", "description for arg1", false, subcmd1Completer)
	subCmd1.AddOption("-arg2", "description for arg2", false, subcmd1Completer)

	subCmd2 := prompter.SubCommand("subcmd2", "subcmd2 description", subcmd2Executor)
	subCmd2.AddOption("-arg1", "description for arg1", false, subcmd2Completer)
	subCmd2.AddOption("-arg2", "description for arg2", false, subcmd2Completer)

	exitCmd := prompter.ExitCommand("exit", "exit the application")

	myCmd1.AddSubCommands(subCmd1, subCmd2)

	comp.RegisterCommands(myCmd1, exitCmd)

	p := prompt.New(comp.Execute, comp.Complete,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("pewpew"),
		prompt.OptionPrefixTextColor(prompt.White),
	)
	p.Run()
}

func subcmd1Executor(args prompter.CmdArgs) error {
	fmt.Printf("subcmd1 args = %s\n", args)
	val1Arg1, err := args.GetFirstValue("-arg1")
	if err != nil {
		return err
	}
	fmt.Println("val1 for arg1", val1Arg1)
	val1Arg2, err := args.GetFirstValue("-arg2")
	if err != nil {
		return err
	}
	fmt.Println("val1 for arg2", val1Arg2)
	return nil
}

func subcmd1Completer(optName string, _ []string) []prompt.Suggest {
	switch optName {
	case "-arg1":
		return []prompt.Suggest{
			{Text: "arg1-option1", Description: "arg1-option1-description"},
			{Text: "arg1-option2", Description: "arg1-option2-description"},
		}
	case "-arg2":
		return []prompt.Suggest{
			{Text: "arg2-option1", Description: "arg2-option1-description"},
			{Text: "arg2-option2", Description: "arg2-option2-description"},
		}
	}
	return []prompt.Suggest{}
}

func subcmd2Executor(args prompter.CmdArgs) error {
	fmt.Printf("subcmd2 args = %s\n", args)
	val1Arg1, err := args.GetFirstValue("-arg1")
	if err != nil {
		return err
	}
	fmt.Println("val1 for arg1", val1Arg1)
	val1Arg2, err := args.GetFirstValue("-arg2")
	if err != nil {
		return err
	}
	fmt.Println("val1 for arg2", val1Arg2)
	return nil
}

func subcmd2Completer(optName string, _ []string) []prompt.Suggest {
	switch optName {
	case "-arg1":
		return []prompt.Suggest{
			{Text: "arg1-option1", Description: "arg1-option1-description"},
			{Text: "arg1-option2", Description: "arg1-option2-description"},
		}
	case "-arg2":
		return []prompt.Suggest{
			{Text: "arg2-option1", Description: "arg2-option1-description"},
			{Text: "arg2-option2", Description: "arg2-option2-description"},
		}
	}
	return []prompt.Suggest{}
}
