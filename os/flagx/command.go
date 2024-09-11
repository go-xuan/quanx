package flagx

import (
	"flag"
)

type Command struct {
	name    string
	usage   string
	options map[string]Option
}

func (cmd *Command) AddOption(options ...Option) *Command {
	for _, option := range options {
		cmd.options[option.Name()] = option
	}
	return cmd
}

func (cmd *Command) Execute(handler func() error) {
	parser.handlers[cmd.name] = handler
}

func (cmd *Command) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	if options := cmd.options; options != nil {
		for _, option := range cmd.options {
			option.Add(fs)
		}
		cmd.options = options
	}
	return fs
}

func AddCommand(name, usage string, options ...Option) *Command {
	if parser == nil {
		parser = &Parser{
			names:    make([]string, 0),
			commands: make(map[string]*Command),
			handlers: make(map[string]func() error),
		}
	}
	if command, ok := parser.commands[name]; ok {
		for _, option := range options {
			command.options[option.Name()] = option
		}
		return command
	} else {
		var opts = make(map[string]Option)
		for _, option := range options {
			opts[option.Name()] = option
		}
		command = &Command{
			name:    name,
			usage:   usage,
			options: opts,
		}
		parser.commands[name] = command
		parser.names = append(parser.names, name)
		return command
	}
}
