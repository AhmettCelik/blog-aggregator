package structure

import "github.com/AhmettCelik/blog-aggregator/internal/config"

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}
