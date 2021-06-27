package main

import (
	"checkport/tcp"
	"fmt"
)

func main() {

	conn, err := tcp.Connect("220.181.38.148", 80)
	if err != nil {
		fmt.Println("ip or port close.")
	}else {
		fmt.Println("ip or port is open.", conn)
	}
}
