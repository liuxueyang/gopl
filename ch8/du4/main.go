package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "enable verbose output")
var sema = make(chan struct{}, 20)
var done = make(chan struct{})

func main() {
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0 {
		roots = []string{"."}
	}

	var wg sync.WaitGroup

	fileSizes := make(chan int64)
	var tick <-chan time.Time

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, fileSizes, &wg)
	}

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	var nfiles, nbytes int64

	for {
		select {
		case <-tick:
			printDiskUsage(nfiles, nbytes)

		case size, ok := <-fileSizes:
			if !ok {
				printDiskUsage(nfiles, nbytes)
				if *verbose {
					fmt.Println("Done")
				}
				return
			} else {
				nfiles++
				nbytes += size
			}
		case <-done:
			for range fileSizes {
			}
			// panic("Cancelled")
		}
	}
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func dirents(dir string) []os.DirEntry {
	select {
	case <-done:
		return nil
	case sema <- struct{}{}:
	}

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

func walkDir(dir string, fileSizes chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done() // NOTE: defer should be called at the start of the function to ensure it runs even if the function exits early

	if cancelled() {
		return
	}

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, fileSizes, wg)
		} else {
			fileInfo, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting info for file %s: %v\n", entry.Name(), err)
				continue
			}
			fileSizes <- fileInfo.Size()
		}
	}
}
