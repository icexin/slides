package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func untar(base string, r io.Reader) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fullpath := filepath.Join(base, hdr.Name)
		info := hdr.FileInfo()
		// as dir
		if info.IsDir() {
			os.MkdirAll(fullpath, 0755)
			continue
		}
		dir := filepath.Dir(fullpath)
		os.MkdirAll(dir, 0755)

		// as file
		f, err := os.Create(fullpath)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, tr)
		if err != nil {
			f.Close()
			return err
		}
		f.Chmod(info.Mode())
		f.Close()
	}
}

func main() {
	untar(".", os.Stdin)
}
