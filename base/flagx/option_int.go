package flagx

import (
	"flag"

	"github.com/go-xuan/quanx/types/anyx"
)

func IntOption(name, usage string, def int) Option {
	return &intOption{
		baseOption: baseOption{
			name:  name,
			usage: usage,
		},
		def: def,
	}
}

type intOption struct {
	baseOption
	def int
}

func (opt *intOption) Name() string {
	return opt.name
}

func (opt *intOption) Usage() string {
	if opt.def != 0 {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *intOption) GetValue(fs *flag.FlagSet) anyx.Value {
	if value := fs.Int(opt.name, opt.def, opt.usage); value != nil && *value != opt.def {
		return anyx.IntValue(*value)
	} else if opt.def != 0 {
		return anyx.IntValue(opt.def)
	}
	return nil
}
