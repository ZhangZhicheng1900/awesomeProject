package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp4", os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(l.Addr().String())
	time.Sleep(100000 * time.Second)
}
