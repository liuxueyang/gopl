package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var use_profile = flag.Bool("profile", false, "Use time command to measure execution")

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
}

func main() {
	flag.Parse()

	src1 := flag.Arg(0)
	src2 := flag.Arg(1)

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: compare-cpp [--time] [--vscode] <file1.cpp> <file2.cpp>")
		os.Exit(1)
	}

	if len(src1) == 0 || len(src2) == 0 {
		fmt.Fprintf(os.Stderr, "Both file paths must be provided.")
		os.Exit(1)
	}

	basename1 := strings.TrimSuffix(filepath.Base(src1), ".cpp")
	basename2 := strings.TrimSuffix(filepath.Base(src2), ".cpp")
	output1 := basename1 + ".output"
	output2 := basename2 + ".output"

	err := compileFile(src1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling the first source file: %v\n", err)
		os.Exit(1)
	}
	err = compileFile(src2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling the second source file: %v\n", err)
		os.Exit(1)
	}

	same, err := compareFiles(output1, output2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error comparing output files: %v\n", err)
		os.Exit(1)
	}

	if same {
		println("Outputs are the same.")
	} else {
		println("Outputs differ.")

		// 检查 vscode 是否可用
		var programName string = "code"
		if !islinux {
			programName = "Code.exe"
		}

		if _, err := exec.LookPath(programName); err != nil {
			fmt.Fprintf(os.Stderr, "%s not found in PATH: %v\n", programName, err)
		} else {
			diffCmd := exec.Command(programName, "--diff", output1, output2)
			diffCmd.Stdin = os.Stdin
			diffCmd.Stdout = os.Stdout
			diffCmd.Stderr = os.Stderr

			err = diffCmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "code execution error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func compileFile(src string) error {
	var options []string
	if *use_profile {
		options = append(options, "--profile")
	}
	options = append(options, "--redirect", src)
	if _, err := exec.LookPath("compile"); err != nil {
		return fmt.Errorf("compile command not found in PATH: %v. Install it from the package ../compile", err)
	}

	return exec.Command("compile", options...).Run()
}

func compareFiles(file1, file2 string) (bool, error) {
	cmd := exec.Command("diff",
		"-Z", // 忽略空格差异
		file1, file2)

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
