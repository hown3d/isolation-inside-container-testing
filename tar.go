package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {
	tr := tar.NewReader(r)

	for {
		hdr, err := tr.Next()
		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case hdr == nil:
			continue
		}

		path := filepath.Join(dst, filepath.Clean(hdr.Name))
		dir := filepath.Dir(path)
		mode := hdr.FileInfo().Mode()
		uid := hdr.Uid
		gid := hdr.Gid

		switch hdr.Typeflag {
		case tar.TypeReg:

			// It's possible a file is in the tar before its directory,
			// or a file was copied over a directory prior to now
			fi, err := os.Stat(dir)
			if os.IsNotExist(err) || !fi.IsDir() {

				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
			}

			currFile, err := os.Create(path)
			if err != nil {
				return err
			}

			if _, err = io.Copy(currFile, tr); err != nil {
				return err
			}

			if err = setFilePermissions(path, mode, uid, gid); err != nil {
				return err
			}

			currFile.Close()
		case tar.TypeDir:
			if err := os.MkdirAll(path, mode); err != nil {
				return err
			}
			if err = setFilePermissions(path, mode, uid, gid); err != nil {
				return err
			}

		case tar.TypeLink:
			// The base directory for a link may not exist before it is created.
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			link := filepath.Clean(filepath.Join(dst, hdr.Linkname))
			if err := os.Link(link, path); err != nil {
				return err
			}

		case tar.TypeSymlink:
			// The base directory for a symlink may not exist before it is created.
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		}
	}
}

func setFilePermissions(path string, mode os.FileMode, uid, gid int) error {
	if err := os.Chown(path, uid, gid); err != nil {
		return err
	}
	// manually set permissions on file, since the default umask (022) will interfere
	// Must chmod after chown because chown resets the file mode.
	return os.Chmod(path, mode)
}
