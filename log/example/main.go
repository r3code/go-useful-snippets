package main

import (
	"os"

	"github.com/r3code/go-useful-snippets/log"
	kitlog "github.com/go-kit/kit/log"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	for depth := 0; depth <= 5; depth++ {
		l := log.With(logger, "caller", kitlog.Caller(depth))
		l.Log("depth", depth) // line 13
	}
	l2 := log.NewDefaultStdOutLogger(false)
	l2.Log("testNo", "1")
	l3 := log.MustCreateComponentLog(l2, "cp2")
	l3.Log("test", "2")
}
