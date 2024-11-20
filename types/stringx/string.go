package stringx

func NewString(t ...string) String {
	var ti = String{notnull: true}
	if len(t) > 0 {
		ti.value = t[0]
	} else {
		ti.value = ""
	}
	return ti
}

type String struct {
	value   string
	notnull bool
}

func (t *String) UnmarshalJSON(bytes []byte) error {
	if l := len(bytes); l >= 0 {
		t.notnull = true
		if l > 1 && bytes[0] == 34 && bytes[l-1] == 34 {
			// 带引号则去掉引号
			t.value = string(bytes[1 : l-1])
		} else {
			// 兼容不带引号的字符串
			t.value = string(bytes)
		}
	} else {
		t.notnull = false
	}
	return nil
}

func (t *String) MarshalJSON() ([]byte, error) {
	if t.notnull {
		var bytes []byte
		bytes = append(bytes, 34)
		bytes = append(bytes, []byte(t.value)...)
		bytes = append(bytes, 34)
		return bytes, nil
	} else {
		return []byte("null"), nil
	}
}

func (t *String) Value() string {
	return t.value
}

func (t *String) NotNull() bool {
	return t.notnull
}
