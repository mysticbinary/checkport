
# 单线程 TCP全连接 端口探测
原理介绍：

TCP全连接端口探测，原理是调用Socket的connect函数连接到目标IP的特定端口上，
如果连接成功说明端口是开放的，如果连接失败，说明端口没有开放。

Golang的net包提供的Dial与DialTimeout函数，对传统的socket函数进行了封装，无论想创建什么协议的连接，
都只需要调用这两个函数即可。这两个函数的区别是DialTimeout增加了超时时间。

调用：
tcp.TcpAllConnect()


# 单线程 TCP半连接 端口探测
TCP半连接端口原理，扫描器只会向目标端口发送一个SYN包，
如果服务器的端口是开放的，会返回SYN/ACK包，如果端口不开放，则会返回RST/ACK包。


调用：
tcp.TcpSynConnect()


# 并发端口扫描流程1 (模版1)

结构理解：
[
    127.0.0.1:[80,81,82],
    127.0.0.2:[80,81,82],
    127.0.0.3:[80,81,82]
]

1. 生成扫描任务列表：
首先解析出需要扫描的IP与端口的切片，
然后将需要扫描的IP与端口列表放入一个 []map[string]int 中，
map的key为IP地址，value为端口，[]map[string]int 表示所有需要扫描的IP与端口对的切片。

2. 分割扫描任务：
根据并发数将需要扫描的 []map[string]int 切片分割成组，
以便按组进行并发扫描。

3. 按组执行扫描任务：
分别将每组扫描任务传入具体的扫描任务中，
扫描任务函数利用sync.WaitGroup实现并发扫描，
在扫描的过程中将结果保存到一个并发安全的map中。

4. 展示扫描结果：
所有扫描任务完成后，
输出保存在并发安全map中的扫描结果。


# 并发端口扫描流程2 (模版2)
toghter1.go对协程的控制不够精细，每组扫描任务都会瞬间启动大量的协程，
然后逐渐关闭，而不是一个平滑的过程。
这种方法可能会瞬间将服务器的CPU占满，为了解决此问题，
在toghter2.go中使用sync.WaitGroup与channel配合实现了新的并发方式.