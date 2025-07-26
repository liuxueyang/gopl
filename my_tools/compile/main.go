package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

var cflags = []string{
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

var (
	colorInfo    = color.Cyan
	colorError   = color.Red
	colorSuccess = color.Green
	colorWarning = color.Yellow
	colorDebug   = color.Magenta
	colorNotice  = color.Blue
)

var profile = flag.Bool("profile", false, "Use time command to measure execution")
var redirect_output = flag.Bool("redirect", false, "Redirect output to file")

var islinux bool

func isLinux() bool {
	uname_output, err := exec.Command("uname").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(uname_output)), "linux")
}

func init() {
	islinux = isLinux()

	home_path, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	debug_path := filepath.Join(home_path, "cpp", "header")
	cflags = append(cflags, "-I", debug_path)
}

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

	var output string

	if *redirect_output {
		output = basename + ".output"
	}

	err := compileFile(src, execname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling source file: %v", err.Error())
		os.Exit(1)
	}

	err = runAndRedirect(execname, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running executable: %v", err.Error())
		os.Exit(1)
	}
}

func compileFile(src string, execname string) error {
	var cmd *exec.Cmd
	srcFile := src

	colorInfo("Compiling source file: %s", srcFile)

	compileCommand := append(cflags, "-o", execname, srcFile)
	cmd = exec.Command("g++", compileCommand...)

	// 捕获标准输出和错误输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed with error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Compiler output:\n%s\n", string(output))
		return err
	}

	if len(output) > 0 {
		colorWarning("%s", string(output))
	}

	return nil
}

func runAndRedirect(execname, outputFile string) (err error) {
	var cmd *exec.Cmd

	colorInfo("Running executable: %s", execname)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if islinux {
		if *profile {
			// 如果使用 -profile 标志，则使用 time 命令来测量执行时间
			cmd = exec.CommandContext(ctx, "/usr/bin/time", "-v", "./"+execname)
		} else {
			cmd = exec.CommandContext(ctx, "./"+execname)
		}
	} else {
		// TODO: time is not supported on Windows
		cmd = exec.CommandContext(ctx, ".\\"+execname)
	}

	if len(outputFile) > 0 {
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

		cmd.Stdout = file
		cmd.Stderr = file_err
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}
