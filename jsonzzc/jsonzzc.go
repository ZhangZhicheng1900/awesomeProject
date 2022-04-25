package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Args struct {
	Command []string `json:"command,omitempty"`
}

type Para struct {
	Action         string `json:"action,omitempty"`
	TimeoutSeconds string `json:"timeoutSeconds,omitempty"`
	Args           Args   `json:"args,omitempty"`
}

type Str struct {
	Str string `json:"string,omitempty"`
}

var x = `/bin/sh&&-c&&abcd.sh`

func main() {
	zzc := Para{}
	zzc.Action = "exec"
	zzc.TimeoutSeconds = "10"
	zzc.Args.Command = strings.Split(x, "&&")
	bytes, _ := json.Marshal(zzc)
	fmt.Println(string(bytes))
	xxx := fmt.Sprintf("%s", string(bytes))
	bbb, _ := json.Marshal(xxx)
	fmt.Println(string(bbb))
}
