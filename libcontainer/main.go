package main

import (
	"log"
	"path/filepath"

	"os"

	"github.com/docker/docker/pkg/reexec"
	"github.com/opencontainers/runc/libcontainer"

	"github.com/hown3d/isolation-tests/internal/image"
	"github.com/opencontainers/runc/libcontainer/specconv"
)

const containerDir = "/var/lib/container"

var rootFs = filepath.Join(containerDir, "alpine")

func main() {
	if reexec.Init() {
		return
	}

	err := image.UnpackImage()
	if err != nil {
		log.Fatal(err)
	}

	runContainer()
}

func runContainer() {
	log.Println("Creating container")
	container, err := createContainer()
	if err != nil {
		log.Fatalf("creating container: %v", err)
	}
	proc := &libcontainer.Process{
		Args:   []string{"/bin/ps"},
		Env:    []string{"PATH=/bin"},
		User:   "root",
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
	specconv.ToRootless(spec)
	spec.Root.Path = rootFs
	config, err := specconv.CreateLibcontainerConfig(&specconv.CreateOpts{
		Spec:            spec,
		RootlessEUID:    os.Geteuid() != 0,
		RootlessCgroups: true,
		CgroupName:      "yeet",
	})
	if err != nil {
		return nil, err
	}
	factory, err := libcontainer.New(containerDir, libcontainer.InitArgs(os.Args[0], "init"))
	if err != nil {
	}
	return factory.Create("container-id", config)
}
