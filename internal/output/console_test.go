package output_test

import (
	"github.com/malyg1n/sql-migrator/internal/output"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	message      = "SomeMessage"
	escCode      = "\u001B"
	successColor = escCode + "[32m"
	infoColor    = escCode + "[34m"
	warningColor = escCode + "[33m"
	errorColor   = escCode + "[31m"
	resetColor   = escCode + "[0m"
)

func TestConsole_TestPrintWithColor(t *testing.T) {
	testCases := []struct {
		name  string
		color string
	}{
		{
			name:  "success",
			color: successColor,
		},
		{
			name:  "error",
			color: errorColor,
		},
		{
			name:  "warning",
			color: warningColor,
		},
		{
			name:  "info",
			color: infoColor,
		},
	}

	console := output.NewConsoleOutput()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			switch tc.name {
			case "success":
				console.PrintSuccess(message)
			case "error":
				console.PrintError(message)
			case "warning":
				console.PrintWarning(message)
			case "info":
				console.PrintInfo(message)
			}

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout
			assert.Equal(t, tc.color+" "+message+"\n"+resetColor, string(out))
		})
	}
}
