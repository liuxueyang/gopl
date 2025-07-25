package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cflags_arr = []string{
	"-std=c++23",
	"-Wall",
	"-Wextra",
	"-DDEBUG",
	"-D_DEBUG",
	"-DLOCAL",
	"-g",
	"-O2",
	"-Wno-unused-result",
}

var use_time = flag.Bool("time", false, "Use time command to measure execution time")

func main() {
	flag.Parse()

	src := flag.Arg(0)
	if flag.NArg() != 1 {
		println("Usage: compile [--time] <file.cpp>")
		return
	}

	if len(src) == 0 {
		println("Source file path must be provided.")
		return
	}

	basename := strings.TrimSuffix(filepath.Base(src), ".cpp")
	var execname string
	if islinux {
		execname = basename + ".out"
	} else {
		execname = basename + ".exe"
	}
	output := basename + ".output"

	err := compileFile(src, execname)
	if err != nil {
		println("Error compiling source file:", err.Error())
		return
	}

	err = runAndRedirect(execname, output)
	if err != nil {
		println("Error running executable:", err.Error())
		return
	}
}

func compileFile(src string, execname string) error {
	if islinux {
		cflags_arr = append(cflags_arr, "-I/home/rakuyo/cpp/header")
	} else {
		// NOTE: Windows 下路径周围不要加双引号
		cflags_arr = append(cflags_arr, "-IC:/Users/rakuy/cpp/header")
	}

	var cmd *exec.Cmd
	srcFile := src

	if islinux {
		compileCommand := append([]string{"g++"}, cflags_arr...)
		compileCommand = append(compileCommand, "-o", execname, srcFile)
		cmd = exec.Command("sh", append([]string{"-c"}, compileCommand...)...)
	} else {
		compileCommand := append(cflags_arr, "-o", execname, srcFile)

		fmt.Printf("Compiling with args: %v\n", compileCommand)
		cmd = exec.Command("g++", compileCommand...)
	}

	// 捕获标准输出和错误输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Compilation failed with error: %v\n", err)
		fmt.Printf("Compiler output:\n%s\n", string(output))
		return err
	}

	if len(output) > 0 {
		fmt.Printf("Compiler output:\n%s\n", string(output))
	}

	return nil
}

func isLinux() bool {
	uname_output, err := exec.Command("uname").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(uname_output)), "linux")
}

var islinux bool

func init() {
	islinux = isLinux()
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

	var cmd *exec.Cmd
	if *use_time && islinux {
		// 如果使用 -time 标志，则使用 time 命令来测量执行时间
		cmd = exec.Command("time", "./"+execname)
		cmd.Stdout = file
		cmd.Stderr = file_err // 可选：也重定向错误输出
	} else {
		// 执行程序并重定向输出
		if islinux {
			cmd = exec.Command("./" + execname)
		} else {
			cmd = exec.Command(".\\" + execname)
		}
		cmd.Stdout = file
		cmd.Stderr = file_err // 可选：也重定向错误输出
	}

	return cmd.Run()
}

// func _debugShowWorkingDir() {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println("Error getting current working directory:", err)
// 		return
// 	}
// 	fmt.Println("Current working directory:", cwd)
// }
