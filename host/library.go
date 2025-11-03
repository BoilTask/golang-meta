package host

import (
	"net"
)

func GetLocalIp() string {
	// 获取本机的所有网络接口
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unkown"
	}
	// 遍历所有接口地址，找到一个非环回的IP地址
	for _, addr := range addrs {
		// 检查地址类型是否为IP地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 返回找到的第一个非环回的IPv4地址
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "unkown"
}
