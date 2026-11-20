package command

import (
	"slices"
	"strings"
)

type generateOptions struct {
	check   bool
	project string
	out     string
}

func argsToOptions(args []string) (generateOptions, bool) {
	var options generateOptions
	var positionals []string
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" {
			positionals = append(positionals, args[index+1:]...)
			break
		}
		if !strings.HasPrefix(arg, "--") {
			positionals = append(positionals, arg)
			continue
		}
		name, value, hasValue := strings.Cut(arg[2:], "=")
		if name == "check" && !hasValue {
			options.check = true
			continue
		}
		if name != "project" && name != "out" {
			return generateOptions{}, false
		}
		if !hasValue {
			index++
			if index == len(args) {
				return generateOptions{}, false
			}
			value = args[index]
		}
		if name == "project" {
			options.project = value
			continue
		}
		options.out = value
	}
	return options, slices.Equal(positionals, []string{"generate"})
}
