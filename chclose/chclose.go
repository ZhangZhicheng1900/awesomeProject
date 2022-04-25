package main

import (
	"fmt"
	"strconv"
)

func main() {
	ch := make(chan struct{})
	close(ch)
	fmt.Println("ok")
	x, err := strconv.ParseInt("21235", 10, 32)
	fmt.Printf("%d, %v", int32(x), err)
}
