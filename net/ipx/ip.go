package ipx

import (
	"net"
	"strconv"
	"strings"
)

// GetLocalIP 获取本地IP
func GetLocalIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ip := ipNet.IP.To4(); ip != nil {
					return ip.String()
				}
			}
		}
	}
	return ""
}

// GetWLANIP 获取当前机器的WLAN IP
func GetWLANIP() string {
	if netInterfaces, err := net.Interfaces(); err == nil {
		for _, netInterface := range netInterfaces {
			if (netInterface.Flags&net.FlagUp) != 0 && netInterface.Name == "WLAN" {
				addrs, _ := netInterface.Addrs()
				for _, address := range addrs {
					if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
						if ip := ipNet.IP.To4(); ip != nil {
							return ip.String()
						}
					}
				}
			}
		}
	}
	return ""
}

// CheckIpMatch 检测IP是否匹配规则
func CheckIpMatch(rules []string, ip string) bool {
	if len(rules) > 0 && ip != "" {
		for _, rule := range rules {
			if len(strings.Split(rule, `-`)) > 1 {
				// rule == "a.b.c.x-a.b.c.y"
				ruleStart, ruleEnd, _ := strings.Cut(rule, `-`)
				prefix, num := SplitIpByLastPoint(ip)
				startPrefix, minNum := SplitIpByLastPoint(ruleStart)
				endPrefix, maxNum := SplitIpByLastPoint(ruleEnd)
				if prefix == startPrefix && num >= minNum && prefix == endPrefix && num <= maxNum {
					return true
				}
			} else {
				ruleNum, ipNum := strings.Split(rule, `.`), strings.Split(ip, `.`)
				switch len(ruleNum) {
				case 1:
					return ruleNum[0] == "*"
				case 2:
					return ipNum[0] == ruleNum[0] && ruleNum[1] == "*"
				case 3:
					return ipNum[0] == ruleNum[0] && ipNum[1] == ruleNum[1] && ruleNum[2] == "*"
				case 4:
					return ipNum[0] == ruleNum[0] && ipNum[1] == ruleNum[1] && ipNum[2] == ruleNum[2] && (ruleNum[3] == "*" || ruleNum[3] == ipNum[3])
				}
			}
		}
	}
	return false
}

// SplitIpByLastPoint 将IP以最后一个.拆分
func SplitIpByLastPoint(ip string) (string, int) {
	if strings.Contains(ip, `.`) {
		i := strings.LastIndex(ip, ".")
		suf, err := strconv.Atoi(ip[i+1:])
		if err != nil {
			return ip[:i], 0
		}
		return ip[:i], suf
	} else {
		return "", 0
	}
}
