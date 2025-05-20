package flagx

import (
	"flag"

	"github.com/go-xuan/quanx/types/anyx"
)

func BoolOption(name, usage string, def bool) Option {
	return &boolOption{
		baseOption: baseOption{name: name, usage: usage},
		def:        def,
		value:      new(bool),
	}
}

type boolOption struct {
	baseOption
	def   bool
	value *bool
}

func (opt *boolOption) Name() string {
	return opt.name
}

func (opt *boolOption) Usage() string {
	if opt.def {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *boolOption) Add(fs *flag.FlagSet) {
	opt.value = fs.Bool(opt.name, opt.def, opt.usage)
}

func (opt *boolOption) GetValue() anyx.Value {
	if opt.value != nil {
		return anyx.BoolValue(*opt.value)
	} else if opt.def {
		return anyx.BoolValue(opt.def)
	}
	return nil
}
