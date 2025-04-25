package osx

import (
	"net"
	"strconv"
	"strings"
)

// GetLocalIP 获取本地IP
func GetLocalIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		return findValidIP(addrs)
	}
	return ""
}

// GetWLANIP 获取当前机器的WLAN IP
func GetWLANIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, item := range interfaces {
		if (item.Flags&net.FlagUp) != 0 && item.Name == "WLAN" {
			if addrs, e := item.Addrs(); e == nil {
				if ip := findValidIP(addrs); ip != "" {
					return ip
				}
			}
		}
	}
	return ""
}

// findValidIP 从地址列表中查找有效的 IPv4 地址
func findValidIP(addrs []net.Addr) string {
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip := ipNet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}

// CheckIpMatch 检测IP是否匹配规则
func CheckIpMatch(rules []string, ip string) bool {
	if len(rules) == 0 || ip == "" {
		return false
	}

	ipParts := strings.Split(ip, ".")
	for _, rule := range rules {
		if strings.Contains(rule, `-`) {
			ruleStart, ruleEnd, found := strings.Cut(rule, `-`)
			if !found {
				continue
			}
			prefix, num := SplitIpByLastPoint(ip)
			startPrefix, minNum := SplitIpByLastPoint(ruleStart)
			endPrefix, maxNum := SplitIpByLastPoint(ruleEnd)
			if prefix == startPrefix && prefix == endPrefix && num >= minNum && num <= maxNum {
				return true
			}
		} else {
			ruleParts := strings.Split(rule, ".")
			match := true
			for i := 0; i < len(ruleParts); i++ {
				if ruleParts[i] != "*" && ruleParts[i] != ipParts[i] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

// SplitIpByLastPoint 将IP以最后一个.拆分
func SplitIpByLastPoint(ip string) (string, int) {
	lastIndex := strings.LastIndex(ip, ".")
	if lastIndex == -1 {
		return "", 0
	}

	prefix := ip[:lastIndex]
	sufStr := ip[lastIndex+1:]
	suf, err := strconv.Atoi(sufStr)
	if err != nil {
		return prefix, 0
	}
	return prefix, suf
}
