package flagx

import (
	"flag"

	"github.com/go-xuan/quanx/types/anyx"
)

func BoolOption(name, usage string, def bool) Option {
	return &boolOption{
		baseOption: baseOption{
			name:  name,
			usage: usage,
		},
		def: def,
	}
}

type boolOption struct {
	baseOption
	def bool
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

func (opt *boolOption) GetValue(fs *flag.FlagSet) anyx.Value {
	if value := fs.Bool(opt.name, opt.def, opt.usage); value != nil && *value != opt.def {
		return anyx.BoolValue(*value)
	} else if opt.def {
		return anyx.BoolValue(opt.def)
	}
	return nil
}
