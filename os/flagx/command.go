package flagx

import (
	"flag"
	"fmt"

	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

type Command struct {
	name     string
	usage    string
	optNames []string
	options  map[string]Option
	handler  func() error
}

func NewCommand(name, usage string, options ...Option) *Command {
	var opts = make(map[string]Option)
	var optNames []string
	for _, option := range options {
		optName := option.Name()
		optNames = append(optNames, optName)
		opts[optName] = option
	}
	return &Command{
		name:     name,
		usage:    usage,
		optNames: optNames,
		options:  opts,
	}
}

// AddOption 添加参数
func (cmd *Command) AddOption(options ...Option) *Command {
	for _, option := range options {
		optName := option.Name()
		if _, ok := cmd.options[optName]; !ok {
			cmd.optNames = append(cmd.optNames, optName)
		}
		cmd.options[optName] = option
	}
	return cmd
}

func (cmd *Command) GetOptionValue(option string) anyx.Value {
	if opt, ok := cmd.options[option]; ok {
		return opt.Get()
	} else {
		fmt.Printf("option [%s] not found\n", option)
		return nil
	}
}

// SetHandler 设置执行器
func (cmd *Command) SetHandler(handler func() error) *Command {
	cmd.handler = handler
	return cmd
}

// Register 命令注册
func (cmd *Command) Register() {
	assertParser()
	var name = cmd.name
	if cmd.handler == nil {
		panic(fmt.Sprintf("command %s has no handler", name))
	}
	if _, ok := parser.commands[name]; ok {
		panic(fmt.Sprintf("command %s has been registered", name))
	}
	parser.names = append(parser.names, name)
	parser.commands[name] = cmd
}

func (cmd *Command) newFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(cmd.name, flag.ExitOnError)
	if options := cmd.options; options != nil {
		for _, option := range cmd.options {
			option.Add(fs)
		}
		cmd.options = options
	}
	return fs
}

func (cmd *Command) Help() {
	fmt.Printf("OPTIONS OF [%s]:\n", fmtx.Cyan.String(cmd.name))
	for _, optName := range cmd.optNames {
		option := cmd.options[optName]
		fmt.Printf("%-50s %s\n", fmtx.Magenta.String("-"+option.Name()), option.Usage())
	}
	fmt.Println()
}
