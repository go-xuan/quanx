package flagx

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/anyx"
)

var _commands *commands

// 初始化
func initCommands() {
	if _commands == nil {
		_commands = &commands{
			names: make([]string, 0),
			child: make(map[string]*Command),
		}
	}
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

// Register 命令注册
func Register(commands ...*Command) {
	initCommands()
	for _, command := range commands {
		var name = command.name
		if command.handler == nil {
			panic(fmt.Sprintf("[%s]命令未设置执行器！", command.name))
		}
		if _, ok := _commands.child[command.name]; ok {
			panic(fmt.Sprintf("[%s]命令未注册", name))
		}
		_commands.names = append(_commands.names, name)
		_commands.child[name] = command
	}
}

// Execute 启动
func Execute() error {
	if _commands == nil {
		return errorx.New("请先注册命令")
	}
	if args := os.Args; len(args) > 1 {
		var name = strings.ToLower(args[1])
		if command, exist := _commands.child[name]; exist {
			if err := command.newFlagSet().Parse(args[2:]); err != nil {
				return errorx.Wrap(err, "解析参数命令失败")
			}
			if handler := command.handler; handler != nil {
				fmtx.Cyan.XPrintf("======当前执行命令是: %s======\n", command.name)
				if err := handler(); err != nil {
					return errorx.Wrap(err, command.name+" handle failed")
				}
			} else {
				fmtx.Red.XPrintf("[%s]命令未设置执行器！", command.name)
			}
			return nil
		}
	}
	_commands.Help()
	return nil
}

func GetOptionValue(command, option string) anyx.Value {
	if cmd, exist := _commands.child[command]; exist {
		if opt, ok := cmd.options[option]; ok {
			return opt.Get()
		} else {
			fmt.Printf("[%s]命令未注册此参数：%s\n", command, option)
			return nil
		}
	} else {
		fmt.Printf("[%s]命令未注册此参数：%s\n", command, option)
		return nil
	}
}

type commands struct {
	names []string
	child map[string]*Command
}

func (p *commands) Help() {
	fmt.Println("\n已注册命令：")
	for _, name := range p.names {
		fmt.Printf("%-50s %s\n", fmtx.Cyan.String(name), p.child[name].usage)
	}
	fmt.Println("\n已注册参数：")
	for _, name := range p.names {
		p.child[name].Help()
	}
}

type Command struct {
	name     string
	usage    string
	optNames []string
	options  map[string]Option
	handler  func() error
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
		fmt.Printf("[%s]参数未找到\n", option)
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
	initCommands()
	var name = cmd.name
	if cmd.handler == nil {
		panic(fmt.Sprintf("[%s]命令未设置执行器！", name))
	}
	if _, ok := _commands.child[name]; ok {
		panic(fmt.Sprintf("[%s]命令未注册！", name))
	}
	_commands.names = append(_commands.names, name)
	_commands.child[name] = cmd
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
	fmt.Printf("[%s]命令的可用参数列表：\n", fmtx.Cyan.String(cmd.name))
	for _, optName := range cmd.optNames {
		option := cmd.options[optName]
		fmt.Printf("%-50s %s\n", fmtx.Magenta.String("-"+option.Name()), option.Usage())
	}
	fmt.Println()
}
