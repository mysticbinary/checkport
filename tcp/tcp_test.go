package tcp

import (
	"checkport/tcp"
	"testing"
)

//go test -v tcp/tcp_test.go tcp/tcp.go
func TestTcpall(t *testing.T) {
	tcp.TcpAllconnect("127.0.0.1", 80)
}

//go test -v -bench . tcp/tcp_test.go tcp/tcp.go
func BenchmarkTcpall(t *testing.B) {
	for i := 0; i < t.N; i++ {
		tcp.TcpAllconnect("127.0.0.1", 80)
		//tcp.TcpSynConnect("127.0.0.1", 80)
	}
}
