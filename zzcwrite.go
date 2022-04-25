package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fileName := os.Args[1]
	writeOffsetK, _ := strconv.Atoi(os.Args[2])
	writeKB, _ := strconv.Atoi(os.Args[3])

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	zzcBytes := make([]byte, writeKB*1024)
	n, errW := f.WriteAt(zzcBytes, int64(writeOffsetK*1024))
	if errW != nil {
		fmt.Println(errW)
		return
	}
	err = f.Sync()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("wrote %dB at offset %dKi in file %s\n", n, writeOffsetK, fileName)

}
