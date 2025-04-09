package flagx

import (
	"flag"
	"fmt"

	"github.com/go-xuan/quanx/types/anyx"
)

type Option interface {
	Name() string
	Usage() string
	Add(fs *flag.FlagSet)
	GetValue() anyx.Value
}

func StringOption(name, usage string, def string) Option {
	return &stringOption{name, usage, def, new(string)}
}

func IntOption(name, usage string, def int) Option {
	return &intOption{name, usage, def, new(int)}
}

func Int64Option(name, usage string, def int64) Option {
	return &int64Option{name, usage, def, new(int64)}
}

func BoolOption(name, usage string, def bool) Option {
	return &boolOption{name, usage, def, new(bool)}
}

func FloatOption(name, usage string, def float64) Option {
	return &floatOption{name, usage, def, new(float64)}
}

type stringOption struct {
	name  string
	usage string
	def   string
	value *string
}

func (opt *stringOption) Name() string {
	return opt.name
}

func (opt *stringOption) Usage() string {
	if opt.def != "" {
		return fmt.Sprintf("%s | default: %s", opt.usage, opt.def)
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

type intOption struct {
	name  string
	usage string
	def   int
	value *int
}

func (opt *intOption) Name() string {
	return opt.name
}

func (opt *intOption) Usage() string {
	if opt.def == 0 {
		return fmt.Sprintf("%s | default: %d", opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *intOption) Add(fs *flag.FlagSet) {
	opt.value = fs.Int(opt.name, opt.def, opt.usage)
}

func (opt *intOption) GetValue() anyx.Value {
	return anyx.IntValue(*opt.value)
}

type int64Option struct {
	name  string
	usage string
	def   int64
	value *int64
}

func (opt *int64Option) Name() string {
	return opt.name
}

func (opt *int64Option) Usage() string {
	if opt.def == 0 {
		return fmt.Sprintf("%s | default: %d", opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *int64Option) Add(fs *flag.FlagSet) {
	opt.value = fs.Int64(opt.name, opt.def, opt.usage)
}

func (opt *int64Option) GetValue() anyx.Value {
	return anyx.Int64Value(*opt.value)
}

type boolOption struct {
	name  string
	usage string
	def   bool
	value *bool
}

func (opt *boolOption) Name() string {
	return opt.name
}

func (opt *boolOption) Usage() string {
	if opt.def {
		return fmt.Sprintf("%s | default: %v", opt.usage, opt.def)
	} else {
		return opt.usage
	}
}

func (opt *boolOption) Add(fs *flag.FlagSet) {
	opt.value = fs.Bool(opt.name, opt.def, opt.usage)
}

func (opt *boolOption) GetValue() anyx.Value {
	return anyx.BoolValue(*opt.value)
}

type floatOption struct {
	name  string
	usage string
	def   float64
	value *float64
}

func (opt *floatOption) Name() string {
	return opt.name
}

func (opt *floatOption) Usage() string {
	if opt.def == float64(0) {
		return fmt.Sprintf("%s | default: %f", opt.usage, opt.def)
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
