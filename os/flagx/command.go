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
			if err := command.Execute(args[2:]); err != nil {
				return errorx.Wrap(err, "命令执行失败")
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
	args     []string          // 当前输入参数
	optNames []string          // 选项名
	options  map[string]Option // 选项列表
	flagSet  *flag.FlagSet     //
	executor func() error      // 命令处理器
}

// Register 命令注册
func (c *Command) Register() {
	var name = c.name
	if c.executor == nil {
		panic(fmt.Sprintf("[%s]命令未设置执行器！", name))
	}
	c.addDefaultOption()
	initManager()
	if _, ok := _manager.commands[name]; ok {
		panic(fmt.Sprintf("[%s]命令未注册！", name))
	}
	_manager.names = append(_manager.names, name)
	_manager.commands[name] = c
}

// Execute 执行
func (c *Command) Execute(args []string) error {
	if err := c.doParse(args); err != nil {
		return errorx.Wrap(err, "命令参数解析失败")
	}
	if err := c.doExecute(); err != nil {
		return errorx.Wrap(err, "命令执行器执行失败")
	}
	return nil
}

func (c *Command) doParse(args []string) error {
	if err := c.GetFlagSet().Parse(args); err != nil {
		return errorx.Wrap(err, "解析命令参数失败")
	}
	c.args = args
	return nil
}

// doExecute 执行
func (c *Command) doExecute() error {
	if executor := c.executor; executor != nil {
		fmtx.Cyan.XPrintf("======当前执行命令是: %s======\n", c.name)
		if err := executor(); err != nil {
			return errorx.Wrap(err, c.name+" execute failed")
		}
	} else {
		fmtx.Red.XPrintf("[%s]命令未设置执行器！", c.name)
	}
	return nil
}

// AddOption 添加参数
func (c *Command) AddOption(options ...Option) *Command {
	for _, option := range options {
		if optName := option.Name(); optName != "" {
			c.options[optName] = option
			if _, ok := c.options[optName]; !ok {
				c.optNames = append(c.optNames, optName)
			}
		}
	}
	return c
}

// GetOptionValue 获取参数值
func (c *Command) GetOptionValue(optName string) anyx.Value {
	if option, ok := c.options[optName]; ok {
		if value := option.GetValue(); value.String() == "-h" {
			_ = c.GetFlagSet().Set("h", "true")
			return anyx.ZeroValue()
		} else {
			return value
		}
	} else {
		fmt.Printf("[%s]参数未找到\n", optName)
		return anyx.ZeroValue()
	}
}

func (c *Command) GetHelpOptionValue() anyx.Value {
	return c.GetOptionValue("h")
}

// SetExecutor 设置执行器
func (c *Command) SetExecutor(executor func() error) *Command {
	c.executor = executor
	return c
}

// GetFlagSet 初始化FlagSet并将参数注册到FlagSet
func (c *Command) GetFlagSet() *flag.FlagSet {
	if c.flagSet == nil {
		fs := flag.NewFlagSet(c.name, flag.ExitOnError)
		if options := c.options; options != nil {
			for _, option := range options {
				option.Add(fs)
			}
			c.options = options
		}
		c.flagSet = fs
	}
	return c.flagSet
}

func (c *Command) GetArgs() []string {
	return c.args
}

func (c *Command) GetArg(index int) string {
	if index > 0 && index < len(c.args) {
		return c.args[index]
	}
	return ""
}

func (c *Command) addDefaultOption() {
	c.AddOption(
		BoolOption("h", "帮助说明", false),
	)
}

// OptionsHelp 命令参数的帮助说明
func (c *Command) OptionsHelp() {
	fmt.Printf("[%s]命令的可用参数列表：\n", fmtx.Cyan.String(c.name))
	for _, optName := range c.optNames {
		option := c.options[optName]
		fmt.Printf("%-50s %s\n", fmtx.Magenta.String("-"+option.Name()), option.Usage())
	}
	fmt.Println()
}
