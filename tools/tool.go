package tools

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func IsRoot() bool {
	return os.Geteuid() == 0
}

func CheckRoot() {
	if !IsRoot() {
		fmt.Println("must run with root")
		os.Exit(0)
	}
}

// 用法示例
//addr, _ := GetIpList("10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24")
func GetIpList(ips string) ([]net.IP, error) {

	addresslist, err := iprange.ParseList(ips)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}
	log.Printf("%+v", addresslist)
	list := addresslist.Expand()
	return list, err
}

// GetPorts(12-20)
func GetPorts(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}
