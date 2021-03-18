package respondents

import "fmt"

const (
	successColor = "\033[32m"
	infoColor    = "\033[34m"
	warningColor = "\033[33m"
	errorColor   = "\033[31m"
	resetColor   = "\033[0m"
)

type ConsoleRespondent struct {
}

func NewConsoleRespondent() *ConsoleRespondent {
	return &ConsoleRespondent{}
}

// Show error in red color
func (r *ConsoleRespondent) PrintError(message string) {
	fmt.Println(errorColor, message)
	fmt.Print(resetColor)
}

// Show success message in green color
func (r *ConsoleRespondent) PrintSuccess(message string) {
	fmt.Println(successColor, message)
	fmt.Print(resetColor)
}

// Show warning message in orange color
func (r *ConsoleRespondent) PrintWarning(message string) {
	fmt.Println(warningColor, message)
	fmt.Print(resetColor)
}

// Show info message in blue color
func (r *ConsoleRespondent) PrintInfo(message string) {
	fmt.Println(infoColor, message)
	fmt.Print(resetColor)
}
