package main

import (
	"log"
	"os"

	"github.com/AhmettCelik/blog-aggregator/internal/commands"
	"github.com/AhmettCelik/blog-aggregator/internal/config"
	"github.com/AhmettCelik/blog-aggregator/internal/structure"
)

func startGator() {
	var cfg config.Config = config.Read()
	var state *structure.State = new(structure.State)
	var cmds *structure.Commands = new(structure.Commands)

	argsWithoutPath := os.Args[1:]
	if len(argsWithoutPath) == 0 {
		log.Fatal("error: not enought arguments were provided.")
	}

	if len(argsWithoutPath) < 2 {
		log.Fatal("error: a username is required.")
	}

	loginCommand := structure.Command{
		Name: argsWithoutPath[0],
		Args: argsWithoutPath,
	}

	state.Config = &cfg

	cmds.Handlers = make(map[string]func(*structure.State, structure.Command) error)
	cmds.Register(loginCommand.Name, commands.HandlerLogin)
	cmds.Run(state, loginCommand)
}
