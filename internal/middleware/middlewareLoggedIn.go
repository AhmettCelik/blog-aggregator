package middleware

import (
	"context"
	"log"
	"os"

	"github.com/AhmettCelik/blog-aggregator/internal/database"
	"github.com/AhmettCelik/blog-aggregator/internal/structure"
)

func MiddlewareLoggedIn(handler func(s *structure.State, cmd structure.Command, user database.User) error) func(*structure.State, structure.Command) error {
	return func(s *structure.State, cmd structure.Command) error {
		user, err := s.Database.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			log.Fatalf("Error getting user from database: %v", err)
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
}
