package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	//	"strings"
)

func getChildPrefix(prefix string, isLast bool) (result string) {
	result = prefix
	if !isLast {
		result += "│"
	}
	result += "\t"
	return result
}

func getPrefix(prefix string, isLast bool) (result string) {
	result = prefix
	if isLast {
		result += "└───"
	} else {
		result += "├───"
	}
	return result
}

func printDirName(out io.Writer, item string, prefix string) (err error) {
	_, err = fmt.Fprintf(out, "%v%v\n", prefix, item)
	if err != nil {
		panic(err)
	}
	return err
}

func getSize(size int) string {
	if size == 0 {
		return "empty"
	}

	return strconv.Itoa(size) + "b"
}

func printFileName(out io.Writer, item string, prefix string, size int64) (err error) {
	_, err = fmt.Fprintf(out, "%v%v (%s)\n", prefix, item, getSize(int(size)))
	if err != nil {
		panic(err)
	}
	return err
}

// By is the type of a "less" function that defines the ordering of its elements arguments.
type By func(p1, p2 *os.FileInfo) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(files []os.FileInfo) {
	fs := &fileSorter{
		files: files,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(fs)
}

// fileSorter joins a By function and a slice of FileInfo to be sorted.
type fileSorter struct {
	files []os.FileInfo
	by    func(p1, p2 *os.FileInfo) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *fileSorter) Len() int {
	return len(s.files)
}

// Swap is part of sort.Interface.
func (s *fileSorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *fileSorter) Less(i, j int) bool {
	return s.by(&s.files[i], &s.files[j])
}

func sortByName(filesInfo []os.FileInfo) (err error) {
	name := func(p1, p2 *os.FileInfo) bool {
		return (*p1).Name() < (*p2).Name()
	}
	// Sort the files by the various criteria.
	By(name).Sort(filesInfo)

	return err
}

func printNodes(out io.Writer, absPath string, areFilesNeed bool, prefix string) error {
	var err error
	fl, err := os.Open(absPath) // For read access.
	if err != nil {
		panic(err)
	}
	fileInfo, err := fl.Stat()
	if err != nil {
		panic(err)
	}
	if !fileInfo.IsDir() {
		//Nothing to do on leaf node
		return nil
	}
	// Max dir size. Memory limitation. Could be problems on merge elements of tree. TODO
	batchSize := 1000000

	// go through all items in dir
	count := 0
	for {
		// get batch of files
		filesInfo, readErr := fl.Readdir(batchSize)
		if readErr != nil && readErr != io.EOF {
			panic(readErr)
		}

		if !areFilesNeed {
			// Filter files in FilesInfo without reallocations
			b := filesInfo[:0]
			count := 0
			for _, x := range filesInfo {
				if x.IsDir() {
					count++
					b = append(b, x)
				}
			}
			filesInfo = filesInfo[:count]
		}
		err = sortByName(filesInfo)
		//sort.sort(items)
		portionCount := len(filesInfo)
		count += portionCount
		// Print all nodes in folder
		for i := 0; i < portionCount; i++ {
			var item = &(filesInfo[i])
			itemPath := absPath + string(os.PathSeparator) + (*item).Name()
			isLast := (i == count-1)
			// Could be optimized by invokes count
			dirPrefix := getPrefix(prefix, isLast)
			// Print item by node type
			if (*item).IsDir() {
				printDirName(out, (*item).Name(), dirPrefix)
				// Process subitems recursively
				printNodes(out, itemPath, areFilesNeed, getChildPrefix(prefix, isLast))
			} else if areFilesNeed {
				printFileName(out, (*item).Name(), dirPrefix, (*item).Size())
			}
		}
		// Does end of directory has been detected
		if readErr == io.EOF {
			break
		}
	}
	return nil
}

func dirTree(out io.Writer, inPath string, areFilesNeed bool) error {
	var err error
	fullPath := inPath
	if !filepath.IsAbs(inPath) {
		fullPath, err = filepath.Abs(inPath)
	}
	if err == nil {
		printNodes(out, fullPath, areFilesNeed, "")
	}

	return err
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
