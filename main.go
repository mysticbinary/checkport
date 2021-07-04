package main

import (
	"checkport/tools"
	"fmt"
)

func main() {

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

	// 测试生成端口
	p, _ := tools.GetPorts("10-20")
	fmt.Println(p)

	// 测试单机并发

	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, err := util.GetIpList(ipList)
		ports, err := util.GetPorts(portList)
		_ = err

		task, _ := scanner.GenerateTask(ips, ports)
		scanner.AssigningTasks(task)
		scanner.PrintResult()

	} else {
		fmt.Printf("%v iplist port\n", os.Args[0])
	}

}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
