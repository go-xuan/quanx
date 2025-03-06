package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	if err := f4(1); err != nil {
		fmt.Println(err)
	}
	if err := f4(2); err != nil {
		fmt.Println(err)
	}
	if err := f4(3); err != nil {
		fmt.Println(err)
	}
	if err := f4(4); err != nil {
		fmt.Println(err)
	}
}

func f1(e int) error {
	switch e {
	case 1:
		return New("携带栈信息的error")
	case 2:
		return errors.New("普通error")
	case 3:
		return fmt.Errorf("fmt.Errorf")
	case 4:
		return fmt.Errorf("fmt.Errorf ==> %w", errors.New("普通error"))
	default:
		return nil
	}
}

func f2(e int) error {
	if err := f1(e); err != nil {
		return err
	}
	return nil
}

func f3(e int) error {
	if err := f2(e); err != nil {
		return err
	}
	return nil
}

func f4(e int) error {
	if err := f3(e); err != nil {
		return err
	}
	return nil
}
