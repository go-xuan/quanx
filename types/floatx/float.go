package floatx

import "strconv"

func NewFloat64(t ...float64) *Float {
	var ti = Float{notnull: true}
	if len(t) > 0 {
		ti.value = t[0]
	} else {
		ti.value = 0
	}
	return &ti
}

type Float struct {
	value   float64
	notnull bool
}

func (t *Float) UnmarshalJSON(bytes []byte) error {
	if value, err := strconv.ParseFloat(string(bytes), 64); err != nil {
		return err
	} else {
		t.value = value
		t.notnull = true
	}
	return nil
}

func (t *Float) MarshalJSON() ([]byte, error) {
	if t.notnull {
		return []byte(strconv.FormatFloat(t.value, 'f', -1, 64)), nil
	} else {
		return []byte("null"), nil
	}
}

func (t *Float) Value() float64 {
	return t.value
}

func (t *Float) NotNull() bool {
	return t.notnull
}
