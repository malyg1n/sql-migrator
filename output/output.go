package output

import "fmt"

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

// Show error in red color
func ShowError(message string) {
	fmt.Println(string(colorRed), message)
	fmt.Print(string(colorReset))
}

// Show success message in green color
func ShowMessage(message string) {
	fmt.Println(string(colorGreen), message)
	fmt.Print(string(colorReset))
}
