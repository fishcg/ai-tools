package logger

import (
	"log"
	"os"
)

var (
	infolog  = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags|log.Lshortfile)
	debuglog = log.New(os.Stdout, "\033[33m[debug]\033[0m ", log.LstdFlags|log.Lshortfile)
	errorlog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	fatallog = log.New(os.Stdout, "\033[31m[fatal]\033[0m ", log.LstdFlags|log.Lshortfile)

	Fatal  = fatallog.Fatal
	Fatalf = fatallog.Fatalf
	Error  = errorlog.Println
	Errorf = errorlog.Printf
	Debug  = debuglog.Println
	Debugf = debuglog.Printf
	Info   = infolog.Println
	Infof  = infolog.Printf
)
