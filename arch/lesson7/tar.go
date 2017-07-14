package main

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

func maketar(dir string, w io.Writer) error {
	base := filepath.Base(dir)
	tr := tar.NewWriter(w)
	defer tr.Close()
	return filepath.Walk(dir, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		p, _ := filepath.Rel(dir, name)
		h.Name = filepath.Join(base, p)
		if err = tr.WriteHeader(h); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			io.Copy(tr, f)
		}
		return nil
	})
}

func main() {
	err := maketar(os.Args[1], os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
