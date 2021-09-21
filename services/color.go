package services

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	ERROR = "[ERROR] "
	WARN  = "[WARNING] "
	INPUT = "[INPUT] "
	INFO  = "[INFO] "
)

func ColorPrint(colorText string, text string, option ...interface{}) {
	if colorText == ERROR {
		fmt.Printf(color.RedString(colorText)+text, option...)
		os.Exit(1)
	} else if colorText == INFO {
		fmt.Printf(color.GreenString(colorText))
	} else if colorText == WARN {
		fmt.Printf(color.YellowString(colorText))
	} else if colorText == INPUT {
		fmt.Printf(color.BlueString(colorText))
	} else {
		fmt.Printf(colorText+text, option...)
		fmt.Println()
		return
	}

	fmt.Printf(text, option...)
	if colorText != INPUT {
		fmt.Println()
	}
}
