package para

import (
	"flag"
	"os"
)

var AllFlags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
