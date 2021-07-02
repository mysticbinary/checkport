package templet

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"checkport/vars"
)

/**
1. 生产扫描任务列表
首先解析出需要扫描的IP与端口的切片，
然后将需要扫描的IP与端口列表放入一个 []map[string]int 中，
map的key为IP地址，value为端口，[]map[string]int 表示所有需要扫描的IP与端口对的切片。

[
    127.0.0.1:[80,81,82],
    127.0.0.2:[80,81,82],
    127.0.0.3:[80,81,82]
]
*/
func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)

	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}

	return tasks, len(tasks)
}

/**
2. 分割扫描任务
根据并发数分割成组，然后将每组任务传入RunTask函数中执行,
len(tasks)%vars.ThreadNum > 0表示len(tasks) / vars.ThreadNum不能整除，
还有剩余的任务列表需要进行处理。
*/
func AssigningTasks(tasks []map[string]int) {
	scanBatch := len(tasks) / vars.ThreadNum

	for i := 0; i < scanBatch; i++ {
		curTask := tasks[vars.ThreadNum*i : vars.ThreadNum*(i+1)]
		RunTask(curTask)
	}

	if len(tasks)%vars.ThreadNum > 0 {
		lastTasks := tasks[vars.ThreadNum*scanBatch:]
		RunTask(lastTasks)
	}
}

/**
3. 按组执行扫描任务，这个版本的并发是通过sync.WaitGroup来控制的，
一次性创建出所有协程，然后等待所有任务完成
*/
func RunTask(tasks []map[string]int) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	// 每次创建len(tasks)个goroutine，每个goroutine只处理一个ip:port对的检测
	for _, task := range tasks {
		for ip, port := range task {
			go func(ip string, port int) {
				err := SaveResult(Connect(ip, port))
				_ = err
				wg.Done()
			}(ip, port)
		}
	}
	wg.Wait()
}

func SaveResult(ip string, port int, err error) error {
	// fmt.Printf("ip: %v, port: %v,err: %v, goruntineNum: %v\n", ip, port, err, runtime.NumGoroutine())
	if err != nil {
		return err
	}

	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)
	}
	return err
}

/**
4. 展示扫描结果，直接通过sync.map的Range方法枚举出所有结果并展示出来
*/
func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("ports: %v\n", value)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}
