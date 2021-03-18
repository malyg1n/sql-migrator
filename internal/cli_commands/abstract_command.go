package cli_commands

import "github.com/malyg1n/sql-migrator/internal"

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
)

// AbstractCommand
type AbstractCommand struct {
	service internal.ServiceContract
}
