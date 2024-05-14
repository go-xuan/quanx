package anyx

import (
	"strconv"
	"strings"
)

type Convert interface {
	String() string
	Int() int
	Int64() int64
	Float64() float64
	Bool() bool
}

func String(v string) Convert {
	return StringConvert{v}
}

type StringConvert struct {
	val string
}

func (c StringConvert) String() string {
	return c.val
}

func (c StringConvert) Int() int {
	if value, err := strconv.Atoi(c.val); err != nil {
		return value
	} else {
		return 0
	}
}

func (c StringConvert) Int64() int64 {
	if value, err := strconv.ParseInt(c.val, 10, 64); err != nil {
		return value
	} else {
		return 0
	}
}

func (c StringConvert) Float64() float64 {
	if value, err := strconv.ParseFloat(c.val, 64); err != nil {
		return value
	} else {
		return 0
	}
}

func (c StringConvert) Bool() bool {
	switch strings.ToLower(c.val) {
	case "true", "是", "yes":
		return true
	default:
		return false
	}
}

type IntConvert struct {
	val int
}

func (c IntConvert) String() string {
	return strconv.Itoa(c.val)
}

func (c IntConvert) Int() int {
	return c.val
}

func (c IntConvert) Int64() int64 {
	return int64(c.val)
}

func (c IntConvert) Float64() float64 {
	return float64(c.val)
}

func (c IntConvert) Bool() bool {
	if c.val == 1 {
		return true
	} else {
		return false
	}
}

type Int64Convert struct {
	val int64
}

func (c Int64Convert) String() string {
	return strconv.FormatInt(c.val, 10)
}

func (c Int64Convert) Int() int {
	return int(c.val)
}

func (c Int64Convert) Int64() int64 {
	return c.val
}

func (c Int64Convert) Float64() float64 {
	return float64(c.val)
}

func (c Int64Convert) Bool() bool {
	if c.val == 1 {
		return true
	} else {
		return false
	}
}

type Float64Convert struct {
	val float64
}

func (c Float64Convert) String() string {
	return strconv.FormatFloat(c.val, 'f', -1, 64)
}

func (c Float64Convert) Int() int {
	return int(c.val)
}

func (c Float64Convert) Int64() int64 {
	return int64(c.val)
}

func (c Float64Convert) Float64() float64 {
	return c.val
}

func (c Float64Convert) Bool() bool {
	if c.val == 1 {
		return true
	} else {
		return false
	}
}
