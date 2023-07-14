package ipx

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 获取本地WLANIP
func GetWLANIP() string {
	return IpMap()["WLAN"]
}

// 获取所有IP
func IpMap() map[string]string {
	var ipMap = make(map[string]string)
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces error:", err.Error())
		return nil
	}
	for _, netInterface := range netInterfaces {
		if (netInterface.Flags & net.FlagUp) != 0 {
			addrs, _ := netInterface.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipMap[netInterface.Name] = ipnet.IP.String()
					}
				}
			}
		}
	}
	return ipMap
}

// 检测IP是否存在
func CheckIpExist(ruleList []string, ip string) bool {
	if ruleList == nil || len(ruleList) == 0 || ip == "" {
		return false
	}
	for _, rule := range ruleList {
		if len(strings.Split(rule, `-`)) > 1 {
			// rule == "a.b.c.x-a.b.c.y"
			beginIp, endIp, _ := strings.Cut(rule, `-`)
			prefix, num := SplitIpByLastPoint(ip)
			prefixBegin, min := SplitIpByLastPoint(beginIp)
			prefixEnd, max := SplitIpByLastPoint(endIp)
			if prefix == prefixBegin && num >= min && prefix == prefixEnd && num <= max {
				return true
			}
		} else {
			ruleNos := strings.Split(rule, `.`)
			ipNos := strings.Split(ip, `.`)
			switch len(ruleNos) {
			case 1:
				return ruleNos[0] == "*"
			case 2:
				if ipNos[0] == ruleNos[0] && ruleNos[1] == "*" {
					return true
				}
			case 3:
				if ipNos[0] == ruleNos[0] && ipNos[1] == ruleNos[1] && ruleNos[2] == "*" {
					return true
				}
			case 4:
				if ipNos[0] == ruleNos[0] && ipNos[1] == ruleNos[1] && ipNos[2] == ruleNos[2] {
					if ruleNos[3] == "*" || ruleNos[3] == ipNos[3] {
						return true
					}
				}
			}
		}
	}
	return false
}

// 将IP以最后一个.拆分
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
