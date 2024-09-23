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
}

func (p *Parser) Help() {
	fmt.Println("ALL COMMANDS：")
	for _, name := range p.names {
		fmt.Printf("%-50s %s\n", fmtx.Cyan.String(name), p.commands[name].usage)
	}
	fmt.Println()
	fmt.Println("ALL OPTIONS：")
	for _, name := range p.names {
		p.commands[name].Help()
	}
}

// 断言
func assertParser() {
	if parser == nil {
		parser = &Parser{
			names:    make([]string, 0),
			commands: make(map[string]*Command),
		}
	}
}

// Register 命令注册
func Register(commands ...*Command) {
	assertParser()
	for _, command := range commands {
		var name = command.name
		if command.handler == nil {
			panic(fmt.Sprintf("command %s has no handler", name))
		}
		if _, ok := parser.commands[command.name]; ok {
			panic(fmt.Sprintf("command %s has been registered", name))
		}
		parser.names = append(parser.names, name)
		parser.commands[name] = command
	}
}

// Execute 启动
func Execute() error {
	if parser == nil {
		return errorx.New("please register the command first")
	}
	if args := os.Args; len(args) > 1 {
		var name = strings.ToLower(args[1])
		if command, exist := parser.commands[name]; exist {
			if err := command.newFlagSet().Parse(args[2:]); err != nil {
				return errorx.Wrap(err, "failed to parse command args")
			}
			if handler := command.handler; handler != nil {
				fmtx.Cyan.XPrintf("======current command is: %s======", command.name)
				fmt.Println()
				if err := handler(); err != nil {
					return errorx.Wrap(err, command.name+" handle failed")
				}
			} else {
				fmtx.Red.XPrintf("command %s hasn't set the handler", command.name)
			}
			return nil
		} else if name == "-h" || name == "-help" {
			parser.Help()
			return nil
		} else {
			return errorx.New(fmt.Sprintf("usage: [app_name] [%s] [options]", strings.Join(parser.names, "|")))
		}
	}
	parser.Help()
	return nil
}

func GetOptionValue(command, option string) anyx.Value {
	if cmd, exist := parser.commands[command]; exist {
		if opt, ok := cmd.options[option]; ok {
			return opt.Get()
		} else {
			fmt.Printf("option [%s] not found\n", option)
			return nil
		}
	} else {
		fmt.Printf("command [%s] not found\n", command)
		return nil
	}
}
