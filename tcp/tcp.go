package tcp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

/**
使用示例
	conn, err := tcp.TcpAllconnect("220.181.38.148", 80)
	if err != nil {
		fmt.Println("ip or port close.")
	} else {
		fmt.Println("ip or port is open.", conn)
	}
TCP全连接
*/
func TcpAllconnect(ip string, port int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 2*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return conn, err
}

// get the local ip and port based on our destination ip
func localIPPort(dstip net.IP) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":54321")
	if err != nil {
		return nil, 0, err
	}
	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port, nil
		}
	}
	return nil, -1, err
}


/**
使用示例
	ip1, port1, err1 := tcp.TcpSynConnect("127.0.0.1", 80)
	if err1 == nil && port1 > 0 {
		fmt.Println("ip or port close.", ip1)
	} else {
		fmt.Printf("ip or port is open.", port1)
	}

TCP半连接
*/
func TcpSynConnect(dstIp string, dstPort int) (string, int, error) {
	srcIp, srcPort, err := localIPPort(net.ParseIP(dstIp))
	dstAddrs, err := net.LookupIP(dstIp)
	if err != nil {
		return dstIp, 0, err
	}

	dstip := dstAddrs[0].To4()
	var dstport layers.TCPPort
	dstport = layers.TCPPort(dstPort)
	srcport := layers.TCPPort(srcPort)

	// Our IP header... not used, but necessary for TCP checksumming.
	ip := &layers.IPv4{
		SrcIP:    srcIp,
		DstIP:    dstip,
		Protocol: layers.IPProtocolTCP,
	}
	// Our TCP header
	tcp := &layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		SYN:     true,
	}
	err = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return dstIp, 0, err
	}

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		return dstIp, 0, err
	}
	defer conn.Close()

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
		return dstIp, 0, err
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(4 * time.Second)); err != nil {
		return dstIp, 0, err
	}

	for {
		b := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			return dstIp, 0, err
		} else if addr.String() == dstip.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcport {
					if tcp.SYN && tcp.ACK {
						// log.Printf("%v:%d is OPEN\n", dstIp, dstport)
						return dstIp, dstPort, err
					} else {
						return dstIp, 0, err
					}
				}
			}
		}
	}
}
