package main

import (
	"errors"
	"log"
	"path/filepath"

	"os"
	"runtime"
	// "syscall"

	"github.com/docker/docker/pkg/reexec"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/opencontainers/runc/libcontainer"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"

	"github.com/opencontainers/runc/libcontainer/specconv"
)

const containerDir = "/var/lib/container"

var rootFs = filepath.Join(containerDir, "alpine")

func main() {
	if reexec.Init() {
		return
	}
	ref, err := name.ParseReference("alpine:latest")
	if err != nil {
		log.Fatal(err)
	}
	img, err := remote.Image(ref)
	if err != nil {
		log.Fatal(err)
	}

	layers, err := img.Layers()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(rootFs, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("creating rootFs dir: %v", err)
		}
	}
	for _, l := range layers {
		r, err := l.Uncompressed()
		if err != nil {
			log.Fatal(err)
		}
		err = Untar(rootFs, r)
		if err != nil {
			log.Fatal("UnTar layer: %v", err)
		}
	}


	runContainer()
	// cmd := reexec.Command("runContainer")
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWNS |
	// 		syscall.CLONE_NEWUTS |
	// 		syscall.CLONE_NEWIPC |
	// 		syscall.CLONE_NEWPID |
	// 		syscall.CLONE_NEWNET |
	// 		syscall.CLONE_NEWUSER,
	// 	UidMappings: []syscall.SysProcIDMap{
	// 		{
	// 			ContainerID: 0,
	// 			HostID:      os.Getuid(),
	// 			Size:        1,
	// 		},
	// 	},
	// 	GidMappings: []syscall.SysProcIDMap{
	// 		{
	// 			ContainerID: 0,
	// 			HostID:      os.Getgid(),
	// 			Size:        1,
	// 		},
	// 	},
	// }
	// if err := cmd.Start(); err != nil {
	// 	log.Fatalf("Error starting the reexec.Command - %v\n", err)
	// }
}

func runContainer() {
	log.Println("Creating container")
	container, err := createContainer()
	if err != nil {
		log.Fatalf("creating container: %v", err)
	}
	proc := &libcontainer.Process{
		Args:   []string{"/bin/ls", "/"},
		Env:    []string{"PATH=/bin"},
		User:   "0",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Init:   true,
	}
	log.Println("Running container")
	err = container.Run(proc)
	if err != nil {
		log.Fatalf("running container: %v", err)
	}
}

func createContainer() (libcontainer.Container, error) {
	spec := specconv.Example()
	spec.Root.Path = rootFs
	factory, err := libcontainer.New(containerDir, libcontainer.InitArgs(os.Args[0], "init"))
	if err != nil {
		return nil, err
	}
	config, err := specconv.CreateLibcontainerConfig(&specconv.CreateOpts{
		Spec:         spec,
		RootlessEUID: os.Geteuid() != 0,
	})
	if err != nil {
		return nil, err
	}
	return factory.Create("container-id", config)
}

func init() {
	// reexec.Register("runContainer", runContainer)
	if len(os.Args) > 1 && os.Args[1] == "init" {
		runtime.GOMAXPROCS(1)
		runtime.LockOSThread()
		factory, _ := libcontainer.New("")
		if err := factory.StartInitialization(); err != nil {
			log.Fatal(err)
		}
		panic("--this line should have never been executed, congratulations--")
	}
}
