package commands

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	godotenv.Load(".env.testing")
	code := m.Run()
	os.Exit(code)
}

