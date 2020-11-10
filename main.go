package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type node struct {
	path     string
	children []node
	empty    bool
}

func main() {
	dryRun := flag.Bool("dry-run", false, "If set, folders that would be deleted will be logged to the console")
	flag.Parse()

	deleter := func(path string) error {
		if *dryRun {
			log.Println(path)
		} else {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
		}

		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory", err)
	}

	rootNode := node{
		path: cwd,
	}

	createTree(&rootNode)

	cleanTree(&rootNode, deleter)
}

func createTree(n *node) {
	files, err := ioutil.ReadDir(n.path)
	if err != nil {
		log.Println("Couldn't read", n.path, "skipping...", err)
		return
	}

	if len(files) == 0 {
		n.empty = true

		return
	}

	n.children = make([]node, 0)

	for _, file := range files {
		childNode := node{
			path: filepath.Join(n.path, file.Name()),
		}

		if file.IsDir() {
			createTree(&childNode)
		}

		n.children = append(n.children, childNode)
	}

	n.empty = true
	for _, child := range n.children {
		if !child.empty {
			n.empty = false
			break
		}
	}
}

func cleanTree(dir *node, deleter func(path string) error) {
	for _, child := range dir.children {
		if child.empty {
			err := deleter(child.path)
			if err != nil {
				log.Println("There was an error deleting", child.path, err)
			}
		} else {
			cleanTree(&child, deleter)
		}
	}
}
