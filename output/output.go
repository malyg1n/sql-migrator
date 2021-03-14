package output

import "fmt"

const (
	successColor = "\033[32m"
	infoColor    = "\033[34m"
	warningColor = "\033[33m"
	errorColor   = "\033[31m"
	resetColor   = "\033[0m"
)

// Show error in red color
func ShowError(message string) {
	fmt.Println(errorColor, message)
	fmt.Print(resetColor)
}

// Show success message in green color
func ShowMessage(message string) {
	fmt.Println(successColor, message)
	fmt.Print(resetColor)
}

// Show warning message in orange color
func ShowWarning(message string) {
	fmt.Println(warningColor, message)
	fmt.Print(resetColor)
}

// Show info message in blue color
func ShowInfo(message string) {
	fmt.Println(infoColor, message)
	fmt.Print(resetColor)
}
