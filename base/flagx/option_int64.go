package flagx

import (
	"flag"
	"github.com/go-xuan/quanx/types/anyx"
)

func Int64Option(name, usage string, def int64) Option {
	return &int64Option{
		baseOption: baseOption{name: name, usage: usage},
		def:        def,
		value:      new(int64),
	}
}

type int64Option struct {
	baseOption
	def   int64
	value *int64
}

func (opt *int64Option) Name() string {
	return opt.name
}

func (opt *int64Option) Usage() string {
	if opt.def == 0 {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *int64Option) Add(fs *flag.FlagSet) {
	opt.value = fs.Int64(opt.name, opt.def, opt.usage)
}

func (opt *int64Option) GetValue() anyx.Value {
	if opt.value != nil {
		return anyx.Int64Value(*opt.value)
	} else if opt.def > 0 {
		return anyx.Int64Value(opt.def)
	}
	return nil
}
