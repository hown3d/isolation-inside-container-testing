package main

import (
	"runtime"
	"github.com/docker/docker/pkg/reexec"
	"os"
	"log"

	_ "github.com/opencontainers/runc/libcontainer/nsenter"

	"github.com/opencontainers/runc/libcontainer"
)

func init() {
	reexec.Register("runContainer", runContainer)
	if len(os.Args) > 1 && os.Args[1] == "init" {
		log.Println("starting initialization of container")
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			log.Fatalf("container initializiation: %v", err)
		}
		panic("--this line should have never been executed, congratulations--")
	}
}
