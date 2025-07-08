package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
			if info.Mode()&0111 != 0 {
				// get absolute path of the file
				absPath, err := filepath.Abs(filepath.Join(path, file.Name()))
				if err != nil {
					return fmt.Errorf("failed to get absolute path for file %s: %w", filepath.Join(path, file.Name()), err)
				}
				if skipFiles[absPath] {
					continue
				}

				// prompt for the user to delete the file
				fmt.Printf("%s\ny: delete the file.\nn: Do not delete the file.\nk: Add the file to allowlist.\n(y/n/k) (n is default): ", filepath.Join(path, file.Name()))
				var response string
				fmt.Scanln(&response)
				switch response {
				case "y", "Y":
					err := os.Remove(filepath.Join(path, file.Name()))
					if err != nil {
						return fmt.Errorf("failed to delete file %s: %w", filepath.Join(path, file.Name()), err)
					}
					fmt.Printf("File %s deleted successfully.\n", filepath.Join(path, file.Name()))
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

var dir_path *string = flag.String("dir", ".", "Directory to search for executable files")
var skipFiles map[string]bool

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

	if n >= 2 && buffer[0] == '#' && buffer[1] == '!' {
		return true, nil // This is a script file
	}

	// Check if the file contains any null bytes
	if !bytes.Contains(buffer[:n], []byte{0}) {
		// check normal script file extensions
		ext := strings.ToLower(filepath.Ext(filePath))
		txtExts := []string{".sh", ".py", ".pl", ".rb", ".js", ".php", ".lua", ".awk"}
		for _, txtExt := range txtExts {
			if txtExt == ext {
				return true, nil // This is a text executable file
			}
		}
	}

	return false, nil
}
