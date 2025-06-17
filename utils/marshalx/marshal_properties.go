package marshalx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/magiconair/properties"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/types/anyx"
)

func Properties() Marshal {
	return propertiesImpl{}
}

type propertiesImpl struct{}

func (p propertiesImpl) Name() string {
	return propertiesMethod
}

func (p propertiesImpl) Marshal(v interface{}) ([]byte, error) {
	var lines []string
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		key := val.Type().Field(i).Tag.Get("properties")
		value := val.Field(i).String()
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return []byte(strings.Join(lines, "\n")), nil
}

func (p propertiesImpl) Unmarshal(data []byte, v interface{}) error {
	if err := anyx.MustStructPointer(v); err != nil {
		return errorx.New("the kind must be struct pointer")
	}
	pp, err := properties.Load(data, properties.UTF8)
	if err != nil {
		return errorx.Wrap(err, "load properties error")
	}
	propertiesSetStructValue(pp, reflect.ValueOf(v).Elem())
	return nil
}

func (p propertiesImpl) Read(path string, v interface{}) error {
	if data, err := readFile(path); err != nil {
		return errorx.Wrap(err, "read file error")
	} else {
		return p.Unmarshal(data, v)
	}
}

func (p propertiesImpl) Write(path string, v interface{}) error {
	if data, err := p.Marshal(v); err != nil {
		return errorx.Wrap(err, "properties marshal error")
	} else if err = filex.WriteFile(path, data); err != nil {
		return errorx.Wrap(err, "write file error")
	}
	return nil
}

// 通过反射为结构体赋值
func propertiesSetStructValue(pp *properties.Properties, value reflect.Value) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if tag := value.Type().Field(i).Tag.Get("properties"); tag != "" {
			propertiesSetFieldValue(pp, tag, field)
		}
	}
}

// 通过反射为字段赋值
func propertiesSetFieldValue(pp *properties.Properties, key string, field reflect.Value) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(pp.GetString(key, ""))
	case reflect.Bool:
		field.SetBool(pp.GetBool(key, false))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(int64(pp.GetInt(key, 0)))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetInt(int64(pp.GetInt(key, 0)))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(pp.GetFloat64(key, 0))
	case reflect.Struct:
		propertiesSetStructValue(pp, field)
	case reflect.Pointer:
		propertiesSetStructValue(pp, field.Elem())
	default:
		fmt.Println("unsupported kind:", field.Kind())
	}
}
