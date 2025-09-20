package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var dir_path *string = flag.String("dir", ".", "Directory to search for executable files")
var skipFiles map[string]bool
var dry_run *bool = flag.Bool("dry-run", false, "Dry run, do not delete files")
var non_interactive *bool = flag.Bool("non-interactive", false, "Non-interactive mode, do not prompt for confirmation")
var sema = make(chan struct{}, 20)

// remove_executable removes executable files from the specified directory and prompts the user for confirmation.

func main() {
	flag.Parse()
	skipFiles = make(map[string]bool)

	// read from ~/.config/clean_exec/kept_files.txt
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Failed to get user config directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(configDir, "clean_exec")
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		fmt.Printf("failed to create config directory %s: %v", configPath, err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(configDir, "clean_exec", "kept_files.txt")
	file, err := os.OpenFile(configFilePath, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)

	if err != nil {
		fmt.Printf("Failed to open kept_files.txt: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		skipFiles[scanner.Text()] = true
	}

	var wg sync.WaitGroup
	filePaths := make(chan string)

	if *dir_path == "" {
		return
	}

	wg.Add(1)
	go my_walk_dir(*dir_path, filePaths, &wg)

	go func() {
		wg.Wait()
		close(filePaths)
	}()

	for full_path := range filePaths {
		if *dry_run {
			fmt.Printf("Dry run: remove file %s\n", full_path)
			continue
		}

		if *non_interactive {
			err := os.Remove(full_path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to delete file %s: %w", full_path, err)
				continue
			}
			fmt.Printf("File %s deleted successfully.\n", full_path)
			continue
		}

		// prompt for the user to delete the file
		fmt.Printf("%s\ny: delete the file.\nn: Do not delete the file.\nk: Add the file to allowlist.\n(y/n/k) (n is default): ", full_path)

		var response string
		fmt.Scanln(&response)

		switch response {
		case "y", "Y":
			err := os.Remove(full_path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to delete file %s: %w", full_path, err)
				continue
			}
			fmt.Printf("File %s deleted successfully.\n", full_path)
		case "k":
			// write the absolute path to a file under the directory ~/.config/clean_exec/
			configDir, err := os.UserConfigDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get user config directory: %w", err)
				continue
			}

			configPath := filepath.Join(configDir, "clean_exec")
			// append the absolute path to kept_files.txt
			configFile := filepath.Join(configPath, "kept_files.txt")
			file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open kept_files.txt: %w", err)
				continue
			}

			absPath, _ := filepath.Abs(full_path)
			_, err = file.WriteString(absPath + "\n")
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to write to kept_files.txt: %w", err)
				continue
			}
			file.Close()

			fmt.Printf("File %s kept and recorded in kept_files.txt.\n", absPath)
		}
	}
}

func isTextExecutable(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 检查Unix/Linux/MacOS脚本文件的shebang标记
	if n >= 2 && buffer[0] == '#' && buffer[1] == '!' {
		return true, nil // This is a script file
	}

	// 检查文件是否包含空字节（二进制文件通常包含空字节）
	if !bytes.Contains(buffer[:n], []byte{0}) {
		// 检查常见脚本文件扩展名
		ext := strings.ToLower(filepath.Ext(filePath))

		// 根据操作系统添加不同的脚本文件扩展名
		var txtExts []string
		if runtime.GOOS == "windows" {
			// Windows平台的脚本文件扩展名
			txtExts = []string{".sh", ".py", ".pl", ".rb", ".js", ".php", ".lua", ".awk", ".ps1", ".bat", ".cmd", ".vbs"}
		} else {
			// Unix/Linux/MacOS平台的脚本文件扩展名
			txtExts = []string{".sh", ".py", ".pl", ".rb", ".js", ".php", ".lua", ".awk"}
		}

		for _, txtExt := range txtExts {
			if txtExt == ext {
				return true, nil // This is a text executable file
			}
		}
	}

	return false, nil
}

func isExecutable(filePath string, info os.FileInfo) bool {
	// 检查文件是否为目录
	if info.IsDir() {
		return false
	}

	// 在Unix/Linux/MacOS系统上检查可执行权限
	if runtime.GOOS != "windows" {
		return info.Mode()&0111 != 0
	} else {
		// 在Windows上，通过文件扩展名判断可执行文件
		ext := strings.ToLower(filepath.Ext(filePath))
		execExts := []string{".exe", ".bat", ".cmd", ".com", ".ps1"}
		for _, execExt := range execExts {
			if ext == execExt {
				return true
			}
		}

		// 也可以通过读取文件头部来判断PE格式（Windows可执行文件格式）
		file, err := os.Open(filePath)
		if err != nil {
			return false
		}
		defer file.Close()

		// 读取文件头部来检查是否为PE格式
		buffer := make([]byte, 2)
		_, err = file.Read(buffer)
		if err != nil {
			return false
		}

		// 检查MZ头（Windows可执行文件的标志）
		return buffer[0] == 'M' && buffer[1] == 'Z'
	}
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

func my_walk_dir(path string, filePaths chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, entry := range dirents(path) {
		if entry.IsDir() {
			wg.Add(1)
			// if file is a directory, recursively call my_walk_dir
			go my_walk_dir(filepath.Join(path, entry.Name()), filePaths, wg)
		} else {
			// if file is a link, skip it
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			full_path := filepath.Join(path, entry.Name())

			// if the file is a script and it is executable, skip it
			isTextExec, err := isTextExecutable(full_path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking if file %s is executable: %v\n", full_path, err)
				continue
			}
			if isTextExec {
				continue
			}

			// if file is not a directory, check if it is executable
			info, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get info for file %s: %w", entry.Name(), err)
				continue
			}

			// check whether the file is executable
			if isExecutable(full_path, info) {
				// get absolute path of the file
				absPath, err := filepath.Abs(full_path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to get absolute path for file %s: %w", full_path, err)
					continue
				}
				if skipFiles[absPath] {
					continue
				}
				filePaths <- full_path
			}
		}
	}
}
