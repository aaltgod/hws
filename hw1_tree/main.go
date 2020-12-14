package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

var isLastChild = make(map[int]bool)

func addSymbols(out io.Writer, dir []os.FileInfo, i, childLvl int, name string) {
	if i == len(dir)-1 {
		io.WriteString(out, "└───"+name+"\n")
		isLastChild[childLvl] = true
	} else {
		io.WriteString(out, "├───"+name+"\n")
		_, exists := isLastChild[childLvl]
		if exists {
			delete(isLastChild, childLvl)
		}
	}
}

func addTabs(out io.Writer, childLvl int) {
	for i := 0; i < childLvl; i++ {
		_, exists := isLastChild[i]
		if exists {
			//fmt.Print(i)
			io.WriteString(out, "\t")
		} else {
			//fmt.Print(i)
			io.WriteString(out, "│\t")
		}
	}
}

func sortDir(dir []os.FileInfo) []os.FileInfo {
	var fileInfoIdx = make(map[int]os.FileInfo)
	var names []string

	for i, v := range dir {
		name := v.Name()
		names = append(names, name)
		fileInfoIdx[i] = v
	}

	sort.Slice(names, func(p, q int) bool {
		return names[p] < names[q]
	})

	for _, fii := range fileInfoIdx {
		for idx, v := range names {
			if fii.Name() == v {
				dir[idx] = fii
				break
			}
			continue
		}
	}
	return dir
}

func readDir(out io.Writer, path string, pf bool, pwd *os.File, childLvl int) error {
	dir, err := pwd.Readdir(-1)
	if err != nil {
		return fmt.Errorf("DIR: %s", err)
	}

	dir = sortDir(dir)

	for i, v := range dir {
		name := v.Name()

		if v.IsDir() {
			newPath := path + string(os.PathSeparator) + name
			pwd, err = os.Open(newPath)
			if err != nil {
				return fmt.Errorf("NEWPATH PWD: %s", err)
			}

			addTabs(out, childLvl)
			addSymbols(out, dir, i, childLvl, name)

			if err = readDir(out, newPath, pf, pwd, childLvl+1); err != nil {
				return fmt.Errorf("CALL readDir: %s", err)
			}
		} else {
			addTabs(out, childLvl)
			addSymbols(out, dir, i, childLvl, name)
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, pf bool) error {
	pwd, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("PWD: %s", err)
	}

	childLvl := 0

	if err = readDir(out, path, pf, pwd, childLvl); err != nil {
		return err
	}

	fmt.Fprintln(out)

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
