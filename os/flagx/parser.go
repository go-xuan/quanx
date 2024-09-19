package flagx

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

var parser *Parser

type Parser struct {
	names    []string
	commands map[string]*Command
	handlers map[string]func() error
}

func (p *Parser) Help() {
	fmt.Println("ALL COMMANDS：")
	for _, name := range p.names {
		fmt.Printf("%-50s %s\n", fmtx.Cyan.String(name), p.commands[name].usage)
	}
	fmt.Println()
	fmt.Println("ALL OPTIONS：")
	for _, name := range p.names {
		fmt.Printf("OPTIONS OF [%s]:\n", fmtx.Cyan.String(name))
		for _, optName := range p.commands[name].optNames {
			option := p.commands[name].options[optName]
			fmt.Printf("%-50s %s\n", fmtx.Magenta.String("-"+option.Name()), option.Usage())
		}
		fmt.Println()
	}
}

// Execute 启动
func Execute() error {
	if parser == nil {
		return errorx.New("please use AddCommand() to add the command first")
	}
	if args := os.Args; len(args) > 1 {
		var name = strings.ToLower(args[1])
		if command, exist := parser.commands[name]; exist {
			if err := command.FlagSet().Parse(args[2:]); err != nil {
				return errorx.Wrap(err, "failed to parse command args")
			}
			if handler, ok := parser.handlers[name]; ok {
				fmtx.Cyan.XPrintf("current command is: %s", command.name)
				if err := handler(); err != nil {
					return errorx.Wrap(err, command.name+" execute failed")
				}
			} else {
				fmtx.Red.XPrintf("current command hasn't set the Executor：%s", command.name)
			}
			return nil
		} else if name == "-h" || name == "-help" {
			parser.Help()
			return nil
		} else {
			return errorx.New(fmt.Sprintf("Usage: [app_name] [%s] [options]", strings.Join(parser.names, "|")))
		}
	}
	parser.Help()
	return nil
}

// GetOptionValue 获取参数值
func GetOptionValue(command, option string) anyx.Value {
	if cmd, exist := parser.commands[command]; exist {
		if opt, ok := cmd.options[option]; ok {
			return opt.Get()
		} else {
			fmt.Printf("Option [%s] not found\n", option)
			return nil
		}
	} else {
		fmt.Printf("Command [%s] not found\n", command)
		return nil
	}
}
