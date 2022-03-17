package logger

import (
	"log"
	"os"

	"github.com/praveennagaraj97/shopee/pkg/color"
)

// Pritty print error log with red color and stops the application.
func ErrorLogFatal(msg interface{}) {
	logger := log.New(os.Stdout, string(color.Red), log.Ldate|log.Ltime|log.Lshortfile)
	logger.Fatal(msg)
}

// Pritty print success log with green color
func PrintLog(msg string, clr color.Color) {
	logger := log.New(os.Stdout, string(clr), log.Ldate|log.Ltime)

	logger.Println(msg + string(color.Reset))
}
