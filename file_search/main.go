package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func main() {
	waitgroup.Add(1)
	go filesearch("/Users/godfreybafana/projects/", "README.md")
	waitgroup.Wait()
	for _, file := range matches {
		fmt.Println("Matched file", file)
	}
}

func filesearch(root string, fn string) {
	fmt.Println("Seraching in", root)
	files, _ := ioutil.ReadDir(root)

	for _, file := range files {
		if strings.Contains(file.Name(), fn) {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			waitgroup.Add(1)
			filesearch(filepath.Join(root, file.Name()), fn)
		}
	}
	waitgroup.Done()
}
