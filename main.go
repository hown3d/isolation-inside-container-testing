package main

import (
	"errors"
	"log"
	"path/filepath"

	"os"

	"github.com/docker/docker/pkg/reexec"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runtime-spec/specs-go"

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
			log.Fatalf("UnTar layer: %v", err)
		}
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
		Spec:         spec,
		RootlessEUID: os.Geteuid() != 0,
		RootlessCgroups: true,
		CgroupName: "yeet",
	})
	if err != nil {
		return nil, err
	}
	factory, err := libcontainer.New(containerDir, libcontainer.InitArgs(os.Args[0], "init"))
	if err != nil {
		return nil, err
	}
	return factory.Create("container-id", config)
}

