package main

import (
	"checkport/tools"
	"fmt"
)

func main() {
	tools.CheckRoot()


	//conn, err := tcp.TcpAllconnect("220.181.38.148", 80)
	//if err != nil {
	//	fmt.Println("ip or port close.")
	//} else {
	//	fmt.Println("ip or port is open.", conn)
	//}
	//
	//ip1, port1, err1 := tcp.TcpSynConnect("127.0.0.1", 80)
	//if err1 == nil && port1 > 0 {
	//	fmt.Println("ip or port close.", ip1)
	//} else {
	//	fmt.Printf("ip or port is open.", port1)
	//}
	p := tools.GetPorts(10-20)
	fmt.Println(p)
}
