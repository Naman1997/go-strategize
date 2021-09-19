package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func ColorPrint(colorText string, text string, option ...interface{}) {

	green := color.New(color.FgGreen).PrintfFunc()
	red := color.New(color.FgRed).PrintfFunc()
	yellow := color.New(color.FgYellow).PrintfFunc()
	blue := color.New(color.FgBlue).PrintfFunc()

	if strings.Contains(colorText, "ERROR") {
		red(colorText)
		fmt.Printf(text, option...)
		os.Exit(1)
	} else if strings.Contains(colorText, "INFO") {
		green(colorText)
	} else if strings.Contains(colorText, "WARNING") {
		yellow(colorText)
	} else if strings.Contains(colorText, "INPUT") {
		blue(colorText)
	} else {
		fmt.Printf(colorText+text, option...)
		fmt.Println()
		return
	}

	fmt.Printf(text, option...)
	fmt.Println()
}
