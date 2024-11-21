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

var _manager *manager

// 初始化管理器
func initManager() {
	if _manager == nil {
		_manager = &manager{
			names:    make([]string, 0),
			commands: make(map[string]*Command),
		}
	}
}

type manager struct {
	names    []string
	commands map[string]*Command
}

func (m *manager) Help() {
	fmt.Println("\n可用命令列表：")
	for _, name := range m.names {
		fmt.Printf("%-50s %s\n", fmtx.Cyan.String(name), m.commands[name].usage)
	}
}

// NewCommand 初始化命令
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
	initManager()
	for _, command := range commands {
		var commandName = command.name
		if command.executor == nil {
			panic(fmt.Sprintf("[%s]命令未设置执行器！", command.name))
		}
		if _, ok := _manager.commands[command.name]; ok {
			panic(fmt.Sprintf("[%s]命令未注册", commandName))
		}
		command.addDefaultOption()
		_manager.names = append(_manager.names, commandName)
		_manager.commands[commandName] = command
	}
}

// Execute 启动
func Execute() error {
	if _manager == nil {
		return errorx.New("请先注册命令")
	}
	if args := os.Args; len(args) > 1 {
		var commandName = strings.ToLower(args[1])
		if command, exist := _manager.commands[commandName]; exist {
			if err := command.GetFlagSet().Parse(args[2:]); err != nil {
				return errorx.Wrap(err, "解析命令参数失败")
			}
			if executor := command.executor; executor != nil {
				fmtx.Cyan.XPrintf("======当前执行命令是: %s======\n", commandName)
				if err := executor(); err != nil {
					return errorx.Wrap(err, commandName+" execute failed")
				}
			} else {
				fmtx.Red.XPrintf("[%s]命令未设置执行器！", commandName)
			}
			return nil
		}
	}
	_manager.Help()
	return nil
}

// GetCommandOptionValue 获取命令参数值
func GetCommandOptionValue(commandName, optionName string) anyx.Value {
	if command, exist := _manager.commands[commandName]; exist {
		return command.GetOptionValue(optionName)
	} else {
		fmt.Printf("[%s]命令未注册！\n", commandName)
		return nil
	}
}

type Command struct {
	name     string            // 命令名
	usage    string            // 命令用法说明
	optNames []string          // 选项参数
	options  map[string]Option // 选项列表
	flagSet  *flag.FlagSet
	executor func() error // 命令处理器
}

// AddOption 添加参数
func (cmd *Command) AddOption(options ...Option) *Command {
	for _, option := range options {
		if optName := option.Name(); optName != "" {
			cmd.options[optName] = option
			if _, ok := cmd.options[optName]; !ok {
				cmd.optNames = append(cmd.optNames, optName)
			}
		}
	}
	return cmd
}

// GetOptionValue 获取参数值
func (cmd *Command) GetOptionValue(optName string) anyx.Value {
	if option, ok := cmd.options[optName]; ok {
		if value := option.GetValue(); value.String() == "-h" {
			_ = cmd.GetFlagSet().Set("h", "true")
			return anyx.ZeroValue()
		} else {
			return value
		}
	} else {
		fmt.Printf("[%s]参数未找到\n", optName)
		return anyx.ZeroValue()
	}
}

func (cmd *Command) GetHelpOptionValue() anyx.Value {
	return cmd.GetOptionValue("h")
}

// SetExecutor 设置执行器
func (cmd *Command) SetExecutor(executor func() error) *Command {
	cmd.executor = executor
	return cmd
}

// Register 命令注册
func (cmd *Command) Register() {
	var name = cmd.name
	if cmd.executor == nil {
		panic(fmt.Sprintf("[%s]命令未设置执行器！", name))
	}
	cmd.addDefaultOption()
	initManager()
	if _, ok := _manager.commands[name]; ok {
		panic(fmt.Sprintf("[%s]命令未注册！", name))
	}
	_manager.names = append(_manager.names, name)
	_manager.commands[name] = cmd
}

// GetFlagSet 初始化FlagSet并将参数注册到FlagSet
func (cmd *Command) GetFlagSet() *flag.FlagSet {
	if cmd.flagSet == nil {
		fs := flag.NewFlagSet(cmd.name, flag.ExitOnError)
		if options := cmd.options; options != nil {
			for _, option := range options {
				option.Add(fs)
			}
			cmd.options = options
		}
		cmd.flagSet = fs
	}
	return cmd.flagSet
}

func (cmd *Command) addDefaultOption() {
	cmd.AddOption(
		BoolOption("h", "帮助说明", false),
	)
}

// OptionsHelp 命令参数的帮助说明
func (cmd *Command) OptionsHelp() {
	fmt.Printf("[%s]命令的可用参数列表：\n", fmtx.Cyan.String(cmd.name))
	for _, optName := range cmd.optNames {
		option := cmd.options[optName]
		fmt.Printf("%-50s %s\n", fmtx.Magenta.String("-"+option.Name()), option.Usage())
	}
	fmt.Println()
}
