package boolx

import "strings"

func NewBool(t ...bool) Bool {
	var ti = Bool{notnull: true}
	if len(t) > 0 && t[0] {
		ti.value = true
	}
	return ti
}

type Bool struct {
	value   bool
	notnull bool
}

func (t *Bool) UnmarshalJSON(bytes []byte) error {
	if value := string(bytes); value != "" {
		t.notnull = true
		if value = strings.ToLower(value); value == "true" || value == "1" {
			t.value = true
		}
	} else {
		t.notnull = false
	}
	return nil
}

func (t *Bool) MarshalJSON() ([]byte, error) {
	if t.notnull && t.value {
		return []byte("true"), nil
	} else if t.notnull {
		return []byte("false"), nil
	} else {
		return []byte("false"), nil
	}
}

func (t *Bool) Value() bool {
	return t.value
}

func (t *Bool) NotNull() bool {
	return t.notnull
}
