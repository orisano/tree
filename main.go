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
)

type Node interface {
	Name() string
	Children() ([]Node, error)
}

type treePrinter struct {
	w        io.Writer
	maxLevel int
}

func (p *treePrinter) Print(node Node) error {
	fmt.Fprintln(p.w, node.Name())
	if err := p.printChildren(node, "", 1); err != nil {
		return err
	}
	return nil
}

func (p *treePrinter) printChildren(node Node, prefix string, level int) error {
	if level == p.maxLevel+1 {
		return nil
	}

	children, err := node.Children()
	if err != nil {
		return fmt.Errorf("children(node=%+v): %w", node, err)
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
		if err := p.printChildren(c, pp, level+1); err != nil {
			return fmt.Errorf("print children: %w", err)
		}
	}
	return nil
}

0
type FileSystemNode struct {
	path       string
	name       string
	isDir      bool
	dirOnly    bool
	showHidden bool
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
		return nil, fmt.Errorf("read dir(path=%v): %w", n.path, err)
	}

	children := make([]Node, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if name == "." || name == ".." {
			continue
		}
		if !n.showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		if n.dirOnly && !entry.IsDir() {
			continue
		}
		children = append(children, &FileSystemNode{path: filepath.Join(n.path, name), name: name, isDir: entry.IsDir(), dirOnly: n.dirOnly, showHidden: n.showHidden})
	}
	return children, nil
}

type TreeOption struct {
	MaxLevel      int
	DirectoryOnly bool
	ShowHidden    bool
}

func PrintDirTree(w io.Writer, dirname string, option TreeOption) error {
	fi, err := os.Stat(dirname)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("is not a directory (path=%v)", dirname)
	}
	printer := &treePrinter{
		w:        w,
		maxLevel: option.MaxLevel,
	}
	root := &FileSystemNode{
		path:       dirname,
		name:       dirname,
		isDir:      true,
		dirOnly:    option.DirectoryOnly,
		showHidden: option.ShowHidden,
	}
	return printer.Print(root)
}

func main() {
	log.SetPrefix("tree: ")
	log.SetFlags(0)

	var option TreeOption

	flag.IntVar(&option.MaxLevel, "L", -1, "Show files and directories up to 'num' levels of depth")
	flag.BoolVar(&option.DirectoryOnly, "d", false, "List directories only")
	flag.BoolVar(&option.ShowHidden, "a", false, "All files are printed")
	flag.Parse()

	dir := "."
	if flag.NArg() == 1 {
		dir = flag.Arg(0)
	}
	if err := PrintDirTree(os.Stdout, dir, option); err != nil {
		log.Fatal(err)
	}
}
