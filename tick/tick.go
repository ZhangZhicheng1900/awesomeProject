package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

var tick = time.Tick(1 * time.Second)

func worker(name string, initWait time.Duration) {
	time.Sleep(initWait)
	for {
		select {
		case <-tick:
			fmt.Printf("%s goroutine %s get tick \n", time.Now().String(), name)
		}
	}
}

func main() {
	fmt.Println(path.Join("/abc", "/abcd/ed"))
	_, err := os.Lstat("tick.go")
	if os.IsNotExist(err) {
		fmt.Println("is not exist")
	}
	fmt.Printf("err %v\n", err)

	go worker("t1", time.Second*3)
	//go worker("t2", time.Millisecond*100)
	//go worker("t3", time.Millisecond*200)
	//go worker("t4", time.Millisecond*300)
	time.Sleep(time.Minute * 3)
}
