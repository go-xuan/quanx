package errorx

import (
	"errors"
	"fmt"
	"strings"
)

func New(s ...string) error {
	if len(s) == 0 {
		return nil
	} else if len(s) == 1 {
		return errors.New(s[0])
	} else {
		return errors.New(strings.Join(s, " "))
	}
}

func Fmt(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
