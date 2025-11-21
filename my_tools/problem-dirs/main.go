package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Dir struct {
	l    int
	r    int
	name string
}

var dirPath *string = flag.String("path", ".", "target directory")

func main() {
	flag.Parse()

	gap := 50
	var dirNames []Dir

	for i := 0; i <= 4000; i += 50 {
		l := i
		r := l + gap - 1
		var d Dir
		d.l, d.r = l, r
		d.name = fmt.Sprintf("%d-%d", l, r)
		dirNames = append(dirNames, d)
	}

	// for _, d := range dirNames {
	// 	fmt.Println(d)
	// }

	entries, err := os.ReadDir(*dirPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if strings.HasSuffix(fileName, ".cpp") {
			re := regexp.MustCompile(`^(\d+)\.`)
			res := re.FindStringSubmatch(fileName)
			if res == nil {
				continue
			}
			if len(res) >= 2 {
				fmt.Println(fileName, res[1])
				num := res[1]
				x, err := strconv.Atoi(num)
				if err != nil {
					log.Fatalf("conver int failed: %v", err)
				}

				d, err := findDir(x, dirNames)
				if err != nil {
					log.Fatalf("can not find specific directory")
				}
				fmt.Println(d)

				newDirName := strings.TrimSuffix(fileName, ".cpp")
				newDir := path.Join(*dirPath, d.name, newDirName)
				_, err = os.Stat(newDir)
				if os.IsNotExist(err) {
					err := os.MkdirAll(newDir, 0o755)
					if err != nil {
						log.Fatalf("failed to make dir: %v", err)
					}
				}

				fmt.Println(newDir)
				fullPath := path.Join(*dirPath, fileName)
				newFullPath := path.Join(newDir, fileName)

				_, err = os.Stat(newFullPath)
				if !os.IsNotExist(err) {
					log.Fatalf("filename already exists!")
				}

				err = os.Rename(fullPath, newFullPath)
				if err != nil {
					log.Fatalf("failed to rename: %v", err)
				}
			}
		}
	}
}

func findDir(num int, dirs []Dir) (Dir, error) {
	for _, d := range dirs {
		if num >= d.l && num <= d.r {
			return d, nil
		}
	}
	return Dir{}, errors.New("not found")
}
