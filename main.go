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

type treePrinter struct {
	w io.Writer
	readDir func(dirname string) ([]os.FileInfo, error)
}

func (p *treePrinter) Print(dirname string) error {
	fmt.Fprintln(p.w, dirname)
	if err := p.printDir(dirname, ""); err != nil {
		return err
	}
	return nil
}

func (p *treePrinter) printDir(dirname, prefix string) error {
	fileInfos, err := p.readDir(dirname)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %v", dirname)
	}

	size := len(fileInfos)
	for i, fi := range fileInfos {
		name := fi.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		hasNext := i != size-1

		pp := prefix
		if hasNext {
			fmt.Fprint(p.w, prefix, "├─── ")
			pp += "│    "
		} else {
			fmt.Fprint(p.w, prefix, "└─── ")
			pp += "     "
		}
		fmt.Fprintln(p.w, name)
		if fi.IsDir() {
			if err := p.printDir(filepath.Join(dirname, name), pp); err != nil {
				return err
			}
		}
	}
	return nil
}

func PrintTree(w io.Writer, dirname string) error {
	return (&treePrinter{w: w, readDir: ioutil.ReadDir}).Print(dirname)
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
