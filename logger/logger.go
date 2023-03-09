package logger

import (
	"fmt"
)

func Error(s string) {
	fmt.Printf(s)
}

func Errorf(s string, v interface{}) {
	fmt.Printf(s, v)
}

func Debug(s string) {
	fmt.Print(s)
}

func Info(s string) {
	fmt.Print(s)
}
