package output

import "fmt"

const (
	escCode      = "\u001B"
	successColor = escCode + "[32m"
	infoColor    = escCode + "[34m"
	warningColor = escCode + "[33m"
	errorColor   = escCode + "[31m"
	resetColor   = escCode + "[0m"
)

// Console is helper for print messages in color
type Console struct {
}

// NewConsoleOutput returns the new instance
func NewConsoleOutput() *Console {
	return &Console{}
}

// PrintError shows error in red color
func (r *Console) PrintError(message string) {
	fmt.Println(errorColor, message)
	fmt.Print(resetColor)
}

// PrintSuccess shows success message in green color
func (r *Console) PrintSuccess(message string) {
	fmt.Println(successColor, message)
	fmt.Print(resetColor)
}

// PrintWarning shows warning message in orange color
func (r *Console) PrintWarning(message string) {
	fmt.Println(warningColor, message)
	fmt.Print(resetColor)
}

// PrintInfo shows info message in blue color
func (r *Console) PrintInfo(message string) {
	fmt.Println(infoColor, message)
	fmt.Print(resetColor)
}
