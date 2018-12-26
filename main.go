package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func PrintTree(w io.Writer, dirname string) error {
	d, err := filepath.Abs(dirname)
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path")
	}
	fmt.Fprintln(w, d)
	if err := printDir(w, d, ""); err != nil {
		return err
	}
	return nil
}

func printDir(w io.Writer, dirname, prefix string) error {
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %v", dirname)
	}

	size := len(fileInfos)
	for i, fi := range fileInfos {
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		p := prefix
		if i == size-1 {
			fmt.Fprint(w, prefix, "└─── ")
			p += "     "
		} else {
			fmt.Fprint(w, prefix, "├─── ")
			p += "│    "
		}
		fmt.Fprintln(w, fi.Name())
		if fi.IsDir() {
			if err := printDir(w, filepath.Join(dirname, fi.Name()), p); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	log.SetPrefix("tree: ")
	log.SetFlags(0)

	flag.Parse()

	dir := flag.Arg(0)
	if err := PrintTree(os.Stdout, dir); err != nil {
		log.Fatal(err)
	}
}
