package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

var isLastChild = make(map[int]bool)

func addFileSize(out io.Writer, file os.FileInfo) {
	fileSize := file.Size()
	if fileSize == 0 {
		io.WriteString(out, " (empty)\n")
	} else {
		io.WriteString(out, " (" + strconv.FormatInt(fileSize, 10) + "b)\n")
	}
}

func addSymbols(out io.Writer, dir []os.FileInfo, i, childLvl int, name string, this os.FileInfo) {
	if i == len(dir)-1 {
		if this.IsDir() {
			io.WriteString(out, "└───"+name+"\n")
		} else {
			io.WriteString(out, "└───"+name)
		}

		isLastChild[childLvl] = true
	} else {
		if this.IsDir() {
			io.WriteString(out, "├───"+name+"\n")
		} else {
			io.WriteString(out, "├───"+name)
		}

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
			io.WriteString(out, "\t")
		} else {
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

	if !pf{
		var newDir []os.FileInfo
		for _, v := range dir {
			if v.IsDir() {
				newDir = append(newDir, v)
			}
		}
		dir = sortDir(newDir)
	} else {
		dir = sortDir(dir)
	}
	for i, v := range dir {
		name := v.Name()

		if v.IsDir() {
			newPath := path + string(os.PathSeparator) + name
			pwd, err = os.Open(newPath)
			if err != nil {
				return fmt.Errorf("NEWPATH PWD: %s", err)
			}

			addTabs(out, childLvl)
			addSymbols(out, dir, i, childLvl, name, v)

			if err = readDir(out, newPath, pf, pwd, childLvl+1); err != nil {
				return fmt.Errorf("CALL readDir: %s", err)
			}
		} else {
			addTabs(out, childLvl)
			addSymbols(out, dir, i, childLvl, name, v)
			addFileSize(out, v)
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
