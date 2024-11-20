package intx

import "strconv"

func NewInt(t ...int) Int {
	var ti = Int{notnull: true}
	if len(t) > 0 {
		ti.value = t[0]
	} else {
		ti.value = 0
	}
	return ti
}

func NewInt64(t ...int64) Int64 {
	var ti = Int64{notNil: true}
	if len(t) > 0 {
		ti.value = t[0]
	} else {
		ti.value = 0
	}
	return ti
}

type Int struct {
	value   int
	notnull bool
}

func (t *Int) UnmarshalJSON(bytes []byte) error {
	if value, err := strconv.Atoi(string(bytes)); err != nil {
		return err
	} else {
		t.value = value
		t.notnull = true
	}
	return nil
}

func (t *Int) MarshalJSON() ([]byte, error) {
	if t.notnull {
		return []byte(strconv.Itoa(t.value)), nil
	} else {
		return []byte("null"), nil
	}
}

func (t *Int) Value() int {
	return t.value
}

func (t *Int) NotNull() bool {
	return t.notnull
}

type Int64 struct {
	value  int64
	notNil bool
}

func (t *Int64) UnmarshalJSON(bytes []byte) error {
	if value, err := strconv.ParseInt(string(bytes), 10, 64); err != nil {
		return err
	} else {
		t.value = value
		t.notNil = true
	}
	return nil
}

func (t *Int64) MarshalJSON() ([]byte, error) {
	if t.notNil {
		return []byte(strconv.FormatInt(t.value, 10)), nil
	} else {
		return []byte("null"), nil
	}
}

func (t *Int64) Value() int64 {
	return t.value
}

func (t *Int64) NotNull() bool {
	return t.notNil
}
