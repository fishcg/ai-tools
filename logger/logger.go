package logger

import (
	"fmt"
)

func Error(s string) {
	fmt.Println(s)
}

func Errorf(s string, v interface{}) {
	fmt.Printf(s, v)
}

func Debug(s string) {
	fmt.Println(s)
}

func Info(s string) {
	fmt.Println(s)
}
