package utils

import (
	"fmt"
	"net"
)

type ipUtil int
var Ip = new(ipUtil)

func (this *ipUtil) GetLocalIP() (ips []string, err error){
	addrs,err := net.InterfaceAddrs()
	if err != nil{
		fmt.Println("get ip arr failed: ",err)
		return
	}
	for _,addr := range addrs{
		if ipnet,ok := addr.(*net.IPNet);ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				ips = append(ips,ipnet.IP.String())
			}
		}
	}
	return
}