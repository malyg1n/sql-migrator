package commands

import (
	"github.com/malyg1n/sqlx-migrator/configs"
)

const (
	timeFormat        = "20060102150405"
	exitStatusSuccess = iota
	exitStatusError
)

// AbstractCommand
type AbstractCommand struct {
	config *configs.DBConfig
}

func ParseArgs(args []string) {

}
