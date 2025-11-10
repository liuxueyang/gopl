package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	verbose = flag.Bool("v", false, "enable verbose output")
	sema    = make(chan struct{}, 20) // limit to 20 concurrent goroutines
)

type Entry struct {
	idx  int
	size int64
}

type RootEntry struct {
	idx    int
	name   string
	nbytes int64
	nfiles int64
}

func main() {
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0 {
		roots = []string{"."}
	}
	rootEntries := make([]RootEntry, len(roots))
	for i, dirName := range roots {
		rootEntries[i] = RootEntry{
			idx:  i,
			name: dirName,
		}
	}

	var wg sync.WaitGroup

	fileSizes := make(chan Entry)
	var tick <-chan time.Time

	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	for _, root := range rootEntries {
		wg.Add(1)
		go walkDir(root, fileSizes, &wg)
	}

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	for {
		select {
		case <-tick:
			printDiskUsageV1(rootEntries)

		case entry, ok := <-fileSizes:
			if !ok {
				printDiskUsageV1(rootEntries)
				if *verbose {
					fmt.Println("Done")
				}
				return
			} else {
				idx := entry.idx
				rootEntries[idx].nfiles++
				rootEntries[idx].nbytes += entry.size
			}
		}
	}
}

func printDiskUsageV1(roots []RootEntry) {
	for _, entry := range roots {
		fmt.Printf("%s:\n", entry.name)
		printDiskUsage(entry.nfiles, entry.nbytes)
	}
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}

	defer func() {
		<-sema
	}()

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory %s: %v\n", dir, err)
		return nil
	}
	return entries
}

func walkDir(dir RootEntry, fileSizes chan<- Entry, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, entry := range dirents(dir.name) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir.name, entry.Name())
			newDir := dir
			newDir.name = subdir
			go walkDir(newDir, fileSizes, wg)
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting info for file %s: %v\n", entry.Name(), err)
				continue
			}
			fileSize := fileInfo.Size()
			// TODO: handle sparse file
			if fileSize > 1024*1024*1024 {
				fmt.Printf("file %s, size: %.1f GB, dir: %s\n", fileInfo.Name(), float64(fileSize)/1e9, dir)
			}
			fileSizes <- Entry{dir.idx, fileSize}
			// fileSizes <- fileInfo.Size()
		}
	}
}
