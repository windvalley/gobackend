package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"go-web-demo/internal/app/apiserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	apiserver.NewApp("apiserver").Run()
}
