package main

import (
	"awesomeProject/readxls"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	x := readxls.ReadXls(os.Args[1])
	b, _ := json.MarshalIndent(x, "", "\t")
	fmt.Println(string(b))
}
