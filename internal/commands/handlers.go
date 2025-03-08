package commands

import (
	"fmt"
	"log"

	"github.com/AhmettCelik/blog-aggregator/internal/structure"
)

type cliCommand struct {
	name     string
	callback func(s *structure.State, cmd structure.Command) error
}

var handlerLogin = func(s *structure.State, cmd structure.Command) error {

	userName := cmd.Args[1]

	if len(userName) == 0 {
		log.Fatal("Error: Write you username for login")
	}

	err := s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("success: user has been set: %s", cmd.Args[1])
	return nil
}

var handlerRegister = func(s *structure.State, cmd structure.Command) error {

	userName := cmd.Args[1]

	if len(userName) == 0 {
		log.Fatal("Error: Write you username for register")
	}

	return nil
}

func GetCommands() map[string]cliCommand {
	commands := make(map[string]cliCommand)

	commands["login"] = cliCommand{
		name:     "login",
		callback: handlerLogin,
	}

	commands["register"] = cliCommand{
		name:     "register",
		callback: handlerRegister,
	}

	return commands
}
