package structure

import (
	"github.com/AhmettCelik/blog-aggregator/internal/config"
	"github.com/AhmettCelik/blog-aggregator/internal/database"
)

type State struct {
	Database *database.Queries
	Config   *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}
