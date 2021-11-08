package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"gobackend/internal/app/apiserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	apiserver.NewApp("gobackend-apiserver").Run()
}
