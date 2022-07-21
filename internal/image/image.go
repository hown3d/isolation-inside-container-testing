package image

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hown3d/isolation-tests/internal/tar"
)

const containerDir = "/var/lib/container"

var rootFs = filepath.Join(containerDir, "alpine")

func UnpackImage() error {
	ref, err := name.ParseReference("alpine:latest")
	if err != nil {
		return fmt.Errorf("parsing reference: %w", err)
	}
	img, err := remote.Image(ref)
	if err != nil {
		return fmt.Errorf("getting remote image: %w", err)
	}

	layers, err := img.Layers()
	if err != nil {
		return fmt.Errorf("getting layers from image: %w", err)
	}

	if err := os.MkdirAll(rootFs, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("creating rootFs dir: %v", err)
		}
	}
	for _, l := range layers {
		r, err := l.Uncompressed()
		if err != nil {
			return fmt.Errorf("getting reader from uncompressed layer: %w", err)
		}
		err = tar.Untar(rootFs, r)
		if err != nil {
			return fmt.Errorf("UnTar layer: %v", err)
		}
	}
	return nil
}
