package flagx

import (
	"flag"
	"github.com/go-xuan/quanx/types/anyx"
)

func StringOption(name, usage string, def string) Option {
	return &stringOption{
		baseOption: baseOption{name: name, usage: usage},
		def:        def,
		value:      new(string),
	}
}

type stringOption struct {
	baseOption
	def   string
	value *string
}

func (opt *stringOption) Name() string {
	return opt.name
}

func (opt *stringOption) Usage() string {
	if opt.def != "" {
		return genUsage(opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *stringOption) Add(fs *flag.FlagSet) {
	opt.value = fs.String(opt.name, opt.def, opt.usage)
}

func (opt *stringOption) GetValue() anyx.Value {
	return anyx.StringValue(*opt.value)
}
