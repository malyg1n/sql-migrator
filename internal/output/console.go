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

type Console struct {
}

func NewConsoleOutput() *Console {
	return &Console{}
}

// Show error in red color
func (r *Console) PrintError(message string) {
	fmt.Println(errorColor, message)
	fmt.Print(resetColor)
}

// Show success message in green color
func (r *Console) PrintSuccess(message string) {
	fmt.Println(successColor, message)
	fmt.Print(resetColor)
}

// Show warning message in orange color
func (r *Console) PrintWarning(message string) {
	fmt.Println(warningColor, message)
	fmt.Print(resetColor)
}

// Show info message in blue color
func (r *Console) PrintInfo(message string) {
	fmt.Println(infoColor, message)
	fmt.Print(resetColor)
}
