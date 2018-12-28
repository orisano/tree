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

type Node interface {
	Name() string
	Children() ([]Node, error)
}

type treePrinter struct {
	w io.Writer
}

func (p *treePrinter) Print(node Node) error {
	fmt.Fprintln(p.w, node.Name())
	if err := p.printChildren(node, ""); err != nil {
		return err
	}
	return nil
}

func (p *treePrinter) printChildren(node Node, prefix string) error {
	children, err := node.Children()
	if err != nil {
		return errors.Wrapf(err, "failed to get children(node=%+v)", node)
	}

	size := len(children)
	for i, c := range children {
		hasNext := i != size-1

		pp := prefix
		if hasNext {
			fmt.Fprint(p.w, prefix, "├─── ")
			pp += "│    "
		} else {
			fmt.Fprint(p.w, prefix, "└─── ")
			pp += "     "
		}

		fmt.Fprintln(p.w, c.Name())
		if err := p.printChildren(c, pp); err != nil {
			return err
		}
	}
	return nil
}

type FileSystemNode struct {
	path  string
	name  string
	isDir bool
}

func (n *FileSystemNode) Name() string {
	return n.name
}

func (n *FileSystemNode) Children() ([]Node, error) {
	if !n.isDir {
		return nil, nil
	}

	entries, err := ioutil.ReadDir(n.path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read dir(path=%v)", n.path)
	}

	children := make([]Node, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		children = append(children, &FileSystemNode{path: filepath.Join(n.path, name), name: name, isDir: entry.IsDir()})
	}
	return children, nil
}

func PrintDirTree(w io.Writer, dirname string) error {
	fi, err := os.Stat(dirname)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.Errorf("is not a directory (path=%v)", dirname)
	}
	return (&treePrinter{w: w}).Print(&FileSystemNode{path: dirname, name: dirname, isDir: true})
}

func main() {
	log.SetPrefix("tree: ")
	log.SetFlags(0)

	flag.Parse()

	dir := flag.Arg(0)
	if err := PrintDirTree(os.Stdout, dir); err != nil {
		log.Fatal(err)
	}
}
