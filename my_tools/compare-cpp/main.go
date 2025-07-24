package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var use_time = flag.Bool("time", false, "Use time command to measure execution time")
var use_vscode_diff = flag.Bool("vscode", false, "Use VSCode diff tool for comparison")

func main() {
	flag.Parse()

	src1 := flag.Arg(0)
	src2 := flag.Arg(1)

	if flag.NArg() != 2 {
		println("Usage: compare-cpp [--time] [--vscode] <file1.cpp> <file2.cpp>")
		return
	}

	if len(src1) == 0 || len(src2) == 0 {
		println("Both file paths must be provided.")
		return
	}

	basename1 := strings.TrimSuffix(filepath.Base(src1), ".cpp")
	basename2 := strings.TrimSuffix(filepath.Base(src2), ".cpp")
	execname1 := basename1 + ".out"
	execname2 := basename2 + ".out"
	output1 := basename1 + ".output"
	output2 := basename2 + ".output"

	err := compileFile(src1, execname1)
	if err != nil {
		println("Error compiling first source file:", err.Error())
		return
	}
	err = compileFile(src2, execname2)
	if err != nil {
		println("Error compiling second source file:", err.Error())
		return
	}

	err = runAndRedirect(execname1, output1)
	if err != nil {
		println("Error executing first compiled file:", err.Error())
		return
	}
	err = runAndRedirect(execname2, output2)
	if err != nil {
		println("Error executing second compiled file:", err.Error())
		return
	}

	same, err := compareFiles(output1, output2)
	if err != nil {
		println("Error comparing output files:", err.Error())
		return
	}

	if same {
		println("Outputs are the same.")
	} else {
		println("Outputs differ.")

		// 检查 vimdiff 是否可用
		if _, err := exec.LookPath("vimdiff"); err != nil {
			println("vimdiff not found in PATH:", err.Error())
		} else {
			println("vimdiff found, executing...")
			diffCmd := exec.Command("vimdiff", output1, output2)
			diffCmd.Stdin = os.Stdin
			diffCmd.Stdout = os.Stdout
			diffCmd.Stderr = os.Stderr
			err = diffCmd.Run()
			if err != nil {
				println("vimdiff execution error:", err.Error())
			}
		}
	}
}

func isLinux() (bool, error) {
	uname_output, err := exec.Command("uname").Output()
	if err != nil {
		println("Error executing uname command:", err.Error())
		return false, err
	}
	return strings.Contains(strings.ToLower(string(uname_output)), "linux"), nil
}

func compileFile(src string, execname string) error {
	cflags := "-std=c++23 -Wall -Wextra -DDEBUG -D_DEBUG -DLOCAL -g -O2 -Wno-unused-result"

	if islinux {
		cflags += " -I/home/rakuyo/cpp/header"
	} else {
		cflags += " -I\"C:\\Users\\rakuy\\cpp\\header\""
	}

	compileCmd := "g++ " + cflags + " -o " + execname + " " + src

	var err error
	if islinux {
		err = exec.Command("sh", "-c", compileCmd).Run()
	} else {
		err = exec.Command("cmd", "/C", compileCmd).Run()
	}
	return err
}

func runAndRedirect(execname, outputFile string) error {
	// 创建输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	file_err, err := os.Create(strings.TrimSuffix(outputFile, ".output") + ".err")
	if err != nil {
		return err
	}
	defer file_err.Close()

	// 执行程序并重定向输出
	cmd := exec.Command("./" + execname)
	cmd.Stdout = file
	cmd.Stderr = file_err // 可选：也重定向错误输出

	return cmd.Run()
}

func compareFiles(file1, file2 string) (bool, error) {
	var cmd *exec.Cmd

	if islinux {
		cmd = exec.Command("diff",
			// "-Z",
			file1, file2)
	} else {
		cmd = exec.Command("fc", "/B", file1, file2) // Windows 使用 fc
	}

	err := cmd.Run()
	if err != nil {
		// diff 返回非零退出码表示文件不同
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return false, nil // 文件不同，但不是错误
			}
		}
		return false, err // 真正的错误
	}

	return true, nil // 文件相同
}

var islinux bool

func init() {
	var err error
	islinux, err = isLinux()
	if err != nil {
		println("Error determining if the system is Linux:", err.Error())
		os.Exit(1)
	}
}
