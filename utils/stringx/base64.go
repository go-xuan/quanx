package stringx

import (
	"bytes"
	"reflect"
)

var Base64Code = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=")

var Base64SafeCode = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_=")

const BitCount = 6

func Encode(data []byte) string {
	return encodeBase64(data, Base64Code)
}

func Decode(str string) []byte {
	return decodeBase64(str, Base64Code)
}

// '+' -> '-',  '/' -> '_'
func EncodeSafeUrl(data []byte) string {
	return encodeBase64(data, Base64SafeCode)
}

// '+' -> '-',  '/' -> '_'
func DecodeSafeUrl(str string) []byte {
	return decodeBase64(str, Base64SafeCode)
}

// 转Base64码
func encodeBase64(data []byte, base64Code []byte) string {
	var buffer []byte
	var lastCount uint = 0
	var lastNum uint = 0
	var num uint = 0
	for _, code := range data {
		num = uint(uint(code>>(8-BitCount+lastCount))|lastNum) & 0x3F
		buffer = append(buffer, base64Code[num])
		lastCount = 8 - BitCount + lastCount
		lastNum = uint(code&(1<<lastCount-1)) << (BitCount - lastCount)
		if lastCount == 6 {
			buffer = append(buffer, base64Code[lastNum])
			lastCount = 0
			lastNum = 0
		}
	}
	if lastCount > 0 {
		buffer = append(buffer, base64Code[lastNum])
	}
	if lastCount == 4 {
		buffer = append(buffer, '=')
	} else if lastCount == 2 {
		buffer = append(buffer, []byte("==")...)
	}
	return string(buffer)
}

// Base64解码
func decodeBase64(str string, base64Code []byte) []byte {
	var buffer []byte
	var lastCount uint8 = 0 //余下的长度
	var lastNum uint8 = 0   //余下的内容
	var num uint8 = 0
	data := []byte(str)
	for _, c := range data {
		code := getCode(c, base64Code)
		if lastCount > 0 {
			leftOff := 6 - (8 - lastCount)
			num = uint8(code>>leftOff) | lastNum

			if code > 63 && lastNum <= 0 {
				break
			}
			buffer = append(buffer, num)

			lastCount = leftOff
			lastNum = uint8(code << (8 - lastCount))
		} else {
			lastCount = 6
			lastNum = uint8(code << 2)
		}
	}
	return buffer
}

func getCode(code byte, base64Code []byte) uint {
	for index, c := range base64Code {
		if c == code {
			return uint(index)
		}
	}
	panic("index error")
}

func StructByReflect(inStructPtr interface{}) string {
	var result bytes.Buffer
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}
	// 遍历结构体
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		// 得到tag中的字段名
		key := t.Tag.Get("json")
		result.WriteString(key)
		if f.Type() == nil {
			result.WriteString("")
		} else {
			result.WriteString(f.String())
		}
	}
	return result.String()
}
