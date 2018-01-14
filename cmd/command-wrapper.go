package cmd

import (
	"errors"
	"strconv"

	"github.com/gi4nks/quant"
)

// -------------------------------
// Cli command wrapper
// -------------------------------
func CmdWrapper(args []string) {
}

func commandWrapper(args []string, cmd quant.Action0) {
	Repository.InitDB()
	Repository.InitSchema()

	CmdWrapper(args)

	cmd()

	defer Repository.CloseDB()
}

// ----------------
// Arguments from command string
// ----------------
func commandFromArguments(args []string) (string, []string, error) {
	if len(args) <= 0 {
		return "", nil, errors.New("Value must be provided!")
	}

	return args[0], Utilities.Tail(args), nil
}

func stringsFromArguments(args []string) ([]string, error) {
	if len(args) <= 0 {
		return nil, errors.New("Value must be provided!")
	}

	return args, nil
}

func stringFromArguments(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Value must be provided!")
	}

	str := args[0]

	return str, nil
}

func intFromArguments(args []string) (int, error) {
	if len(args) != 1 {
		return -1, errors.New("Value must be provided!")
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		return -1, err
	}

	return i, nil
}
