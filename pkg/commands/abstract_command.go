package commands

import (
	"github.com/malyg1n/sql-migrator/pkg/services"
)

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
)

// AbstractCommand
type AbstractCommand struct {
	service services.ServiceContract
}
