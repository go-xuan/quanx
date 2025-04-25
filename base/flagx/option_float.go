package flagx

import (
	"flag"
	"github.com/go-xuan/quanx/types/anyx"
)

func FloatOption(name, usage string, def float64) Option {
	return &floatOption{
		baseOption: baseOption{name: name, usage: usage},
		def:        def,
		value:      new(float64),
	}
}

type floatOption struct {
	baseOption
	def   float64
	value *float64
}

func (opt *floatOption) Name() string {
	return opt.name
}

func (opt *floatOption) Usage() string {
	if opt.def == float64(0) {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *floatOption) Add(fs *flag.FlagSet) {
	opt.value = fs.Float64(opt.name, opt.def, opt.usage)
}

func (opt *floatOption) GetValue() anyx.Value {
	return anyx.Float64Value(*opt.value)
}
