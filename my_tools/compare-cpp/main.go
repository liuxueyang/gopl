package main

import (
	"flag"
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

	islinux, err := isLinux()
	if err != nil {
		println("Error determining if the system is Linux:", err.Error())
		return
	}

	if islinux {
		if *use_time {
			err = exec.Command("sh", "-c", "time ./"+execname1+" > "+output1+" 2> "+basename1+".time").Run()
			if err != nil {
				println("Error executing first program:", err.Error())
				return
			}

			err = exec.Command("sh", "-c", "time ./"+execname2+" > "+output2+" 2> "+basename2+".time").Run()
			if err != nil {
				println("Error executing second program:", err.Error())
				return
			}
		} else {
			err = exec.Command("./" + execname1 + " > " + output1).Run()
			if err != nil {
				println("Error executing first program:", err.Error())
				return
			}
			err = exec.Command("./" + execname2 + " > " + output2).Run()
			if err != nil {
				println("Error executing second program:", err.Error())
				return
			}
		}
	} else {
		if *use_time {
			err = exec.Command("cmd", "/C", "time /T && "+execname1+" > "+output1).Run()
			if err != nil {
				println("Error executing first program:", err.Error())
				return
			}

			err = exec.Command("cmd", "/C", "time /T && "+execname2+" > "+output2).Run()
			if err != nil {
				println("Error executing second program:", err.Error())
				return
			}
		} else {
			err = exec.Command(execname1 + " > " + output1).Run()
			if err != nil {
				println("Error executing first program:", err.Error())
				return
			}

			err = exec.Command(execname2 + " > " + output2).Run()
			if err != nil {
				println("Error executing second program:", err.Error())
				return
			}
		}
	}

	// Compare outputs
	// TODO:
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

	islinux, err := isLinux()
	if err != nil {
		println("Error determining if the system is Linux:", err.Error())
		return err
	}

	if islinux {
		cflags += " -I/home/rakuyo/cpp/header"
	} else {
		cflags += " -I\"C:\\Users\\rakuy\\cpp\\header\""
	}

	compileCmd := "g++ " + cflags + " -o " + execname + " " + src

	if islinux {
		err = exec.Command("sh", "-c", compileCmd).Run()
	} else {
		err = exec.Command("cmd", "/C", compileCmd).Run()
	}
	return err
}
