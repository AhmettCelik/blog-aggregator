package commands

import (
	"fmt"
	"log"

	"github.com/AhmettCelik/blog-aggregator/internal/structure"
)

func HandlerLogin(s *structure.State, cmd structure.Command) error {

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

func HandlerRegister(s *structure.State, cmd structure.Command) error {

	userName := cmd.Args[1]

	if len(userName) == 0 {
		log.Fatal("Error: Write you username for register")
	}

	return nil
}
