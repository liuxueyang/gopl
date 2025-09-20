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
)

var dir_path *string = flag.String("dir", ".", "Directory to search for executable files")
var skipFiles map[string]bool
var dry_run *bool = flag.Bool("dry-run", false, "Dry run, do not delete files")
var non_interactive *bool = flag.Bool("non-interactive", false, "Non-interactive mode, do not prompt for confirmation")

// remove_executable removes executable files from the specified directory and prompts the user for confirmation.

func main() {
	flag.Parse()
	skipFiles = make(map[string]bool)

	// read from ~/.config/rm_exec/kept_files.txt
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Failed to get user config directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(configDir, "rm_exec")
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		fmt.Printf("failed to create config directory %s: %v", configPath, err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(configDir, "rm_exec", "kept_files.txt")
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

	my_walk_dir(*dir_path)
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
func remove_executable(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// if file is executable, add to the list
			if info.Mode()&0111 != 0 {
				files = append(files, filePath)

				// prompt for the user to delete the file
				fmt.Printf("%s\ny: delete the file.\nn: Do not delete the file.\nk: Add the file to allowlist.\n(y/n/k): ", filePath)
				var response string
				fmt.Scanln(&response)
				if response == "y" || response == "Y" {
					err := os.Remove(filePath)
					if err != nil {
						return fmt.Errorf("failed to delete file %s: %w", filePath, err)
					}
					fmt.Printf("File %s deleted successfully.\n", filePath)
				} else {
					fmt.Printf("File %s not deleted.\n", filePath)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func my_walk_dir(path string) error {
	// list all files in the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			// if file is a directory, recursively call my_walk_dir
			err := my_walk_dir(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		} else {
			// if file is a link, skip it
			if file.Type()&os.ModeSymlink != 0 {
				continue
			}

			// if the file is a script and it is executable, skip it
			isTextExec, err := isTextExecutable(filepath.Join(path, file.Name()))
			if err != nil {
				fmt.Printf("Error checking if file %s is executable: %v\n", filepath.Join(path, file.Name()), err)
				continue
			}
			if isTextExec {
				continue
			}

			// if file is not a directory, check if it is executable
			info, err := file.Info()
			if err != nil {
				return fmt.Errorf("failed to get info for file %s: %w", file.Name(), err)
			}

			full_path := filepath.Join(path, file.Name())
			// check whether the file is executable
			if isExecutable(full_path, info) {
				// get absolute path of the file
				absPath, err := filepath.Abs(full_path)
				if err != nil {
					return fmt.Errorf("failed to get absolute path for file %s: %w", full_path, err)
				}
				if skipFiles[absPath] {
					continue
				}

				if *dry_run {
					fmt.Printf("Dry run: remove file %s", full_path)
					continue
				}

				if *non_interactive {
					err := os.Remove(full_path)
					if err != nil {
						return fmt.Errorf("failed to delete file %s: %w", full_path, err)
					}
					fmt.Printf("File %s deleted successfully.\n", full_path)
					continue
				}

				// prompt for the user to delete the file
				fmt.Printf("%s\ny: delete the file.\nn: Do not delete the file.\nk: Add the file to allowlist.\n(y/n/k) (n is default): ", filepath.Join(path, file.Name()))
				var response string
				fmt.Scanln(&response)
				switch response {
				case "y", "Y":
					err := os.Remove(full_path)
					if err != nil {
						return fmt.Errorf("failed to delete file %s: %w", full_path, err)
					}
					fmt.Printf("File %s deleted successfully.\n", full_path)
				case "k":
					// write the absolute path to a file under the directory ~/.config/rm_exec/
					configDir, err := os.UserConfigDir()
					if err != nil {
						return fmt.Errorf("failed to get user config directory: %w", err)
					}

					configPath := filepath.Join(configDir, "rm_exec")
					// append the absolute path to kept_files.txt
					configFile := filepath.Join(configPath, "kept_files.txt")
					file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return fmt.Errorf("failed to open kept_files.txt: %w", err)
					}
					defer file.Close()

					_, err = file.WriteString(absPath + "\n")
					if err != nil {
						return fmt.Errorf("failed to write to kept_files.txt: %w", err)
					}

					fmt.Printf("File %s kept and recorded in kept_files.txt.\n", absPath)
				}
			}
		}
	}
	return nil
}
