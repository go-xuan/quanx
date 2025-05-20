package flagx

import (
	"flag"

	"github.com/go-xuan/quanx/types/anyx"
)

func IntOption(name, usage string, def int) Option {
	return &intOption{
		baseOption: baseOption{name: name, usage: usage},
		def:        def,
		value:      new(int),
	}
}

type intOption struct {
	baseOption
	def   int
	value *int
}

func (opt *intOption) Name() string {
	return opt.name
}

func (opt *intOption) Usage() string {
	if opt.def == 0 {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *intOption) Add(fs *flag.FlagSet) {
	opt.value = fs.Int(opt.name, opt.def, opt.usage)
}

func (opt *intOption) GetValue() anyx.Value {
	if opt.value != nil {
		return anyx.IntValue(*opt.value)
	} else if opt.def > 0 {
		return anyx.IntValue(opt.def)
	}
	return nil
}
