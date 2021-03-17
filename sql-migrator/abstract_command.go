package sql_migrator

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
)

// AbstractCommand
type AbstractCommand struct {
	service ServiceContract
}
