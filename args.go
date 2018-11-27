package prompter

import "fmt"

// CmdArgs contains the arguments passed to the command.
type CmdArgs map[string][]string

// TODO: Write tests.

// GetValue returns the n-th value of an argument.
// TODO: Remove the error message here and just return "" instead of errors?
func (c CmdArgs) GetValue(key string, n int) (string, error) {
	if !c.Contains(key) {
		return "", fmt.Errorf("%s is not present", key)
	}
	if len(c[key]) == 0 {
		return "", fmt.Errorf("%s has no associated value", key)
	}
	if len(c[key]) < n {
		return "", fmt.Errorf("%s has %d values, got %d", key, len(c[key]), n)
	}
	return c[key][n-1], nil
}

// GetFirstValue returns the first value of an argument.
func (c CmdArgs) GetFirstValue(key string) (string, error) {
	return c.GetValue(key, 1)
}

// Contains returns true if the key exists in args.
func (c CmdArgs) Contains(key string) bool {
	_, exists := c[key]
	return exists
}
