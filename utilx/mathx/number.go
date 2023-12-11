package mathx

import (
	"math"
	"strings"
)

var numbers = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
var units = []string{"个", "十", "百", "千", "万", "十", "百", "千", "亿", "十", "百", "千"}

// 中文转数组映射
var cn2NumMap = InitCn2NumMap()

func InitCn2NumMap() map[string]int {
	var m = make(map[string]int)
	for i, s := range numbers {
		m[s] = i
	}
	return m
}

// 单位级别映射
var powerLevelMap = map[string]*PowerLevel{
	"亿": {int(math.Pow10(8)), true},
	"万": {int(math.Pow10(4)), true},
	"千": {int(math.Pow10(3)), false},
	"百": {int(math.Pow10(2)), false},
	"十": {int(math.Pow10(1)), false},
	"零": {0, false},
}

// 位数级别
type PowerLevel struct {
	pow     int  // 量级
	isLarge bool // 是否大数额
}

// 中文数字转阿拉伯数字
// 0 <= result < 10^12
func ConvertChinese2Arab(input string) (result int) {
	cnRunes := []rune(input)
	cur := 0 // 当前数字
	sum := 0 // 当前总和，当跨亿级或万级时需要清零重新求和
	for i := 0; i < len(cnRunes); i++ {
		item := string(cnRunes[i])
		// 将中文转为阿拉伯数字
		itemNum := cn2NumMap[item]
		if itemNum > 0 { // 当前值：一/二/三/四/五/六/七/八/九
			if cur == 0 {
				cur = itemNum
			} else {
				cur *= 10
				cur += itemNum
			}
			if i == len(cnRunes)-1 { // 汇总计算
				sum += cur
				result += sum
				break
			}
		} else { // 当前值：零/十/百/千/万/亿
			level := powerLevelMap[item]
			if level.isLarge { // 当前值：万/亿
				sum = (sum + cur) * level.pow
				result = result + sum
				sum = 0 // sum归零，重新计算下一个区间
			} else { // 当前值：零/十/百/千
				if cur == 0 {
					sum += level.pow
				} else {
					sum += cur * level.pow
				}
			}
			cur = 0                  // cur归零
			if i == len(cnRunes)-1 { // 汇总计算
				result += sum
				break
			}
		}
	}
	return
}

// 阿拉伯转中文
// 0 <= input < 10^12
func ConvertArabToChinese(input int) (result string) {
	var arabs []int
	for ; input > 0; input = input / 10 {
		// arabs是倒序
		arabs = append(arabs, input%10)
	}
	for i := len(arabs) - 1; i >= 0; i-- {
		result += numbers[arabs[i]] + units[i]
	}
	for {
		temp := result
		result = strings.Replace(result, "零亿", "亿", 1)
		result = strings.Replace(result, "零万", "万", 1)
		result = strings.Replace(result, "零千", "零", 1)
		result = strings.Replace(result, "零百", "零", 1)
		result = strings.Replace(result, "零十", "零", 1)
		result = strings.Replace(result, "零零", "零", 1)
		if result == temp {
			break
		}
	}
	result = strings.Replace(result, "亿万", "亿", 1)
	result = strings.TrimSuffix(result, "零个")
	return
}
