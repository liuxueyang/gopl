package tempconv

import (
	"flag"
	"fmt"
	"gitee.com/liuxueyang/gopl/ch2/tempconv0"
)

// *celsiusFlag 上的方法包括 *celsiusFlag 和 celsiusFlag 的方法
type celsiusFlag struct {
	tempconv0.Celsius
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C":
		f.Celsius = tempconv0.Celsius(value)
		return nil
	case "F":
		f.Celsius = tempconv0.FToC(tempconv0.Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

func CelsiusFlag(name string, value tempconv0.Celsius, usage string) *tempconv0.Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
