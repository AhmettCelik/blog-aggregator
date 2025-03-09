package commands

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AhmettCelik/blog-aggregator/internal/database"
	"github.com/AhmettCelik/blog-aggregator/internal/rss"
	"github.com/AhmettCelik/blog-aggregator/internal/structure"
	"github.com/google/uuid"
)

func HandlerLogin(s *structure.State, cmd structure.Command) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error: a username is required.")
	}

	userName := cmd.Args[1]

	user, err := s.Database.GetUser(context.Background(), userName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Error checking existing user %v", err)
	} else if err == sql.ErrNoRows {
		log.Fatalf("Error: register for login! This username does not exists on our database %s", user.Name)
		os.Exit(1)
	}

	if len(userName) == 0 {
		log.Fatal("Error: Write you username for login")
	}

	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("success: user has been set: %s", cmd.Args[1])
	return nil
}

func HandlerRegister(s *structure.State, cmd structure.Command) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error: You must provide a username for registration")
	}

	now := time.Now()
	userName := cmd.Args[1]

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      userName,
	}

	existingUser, err := s.Database.GetUser(context.Background(), userName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Error checking existing user: %v", err)
	}

	if existingUser.Name == userName {
		log.Fatalf("Error: Username '%s' is already taken.", userName)
		os.Exit(1)
	}

	user, err := s.Database.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Error creating user: %v", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user: %v", err)
	}

	userJsonFormat, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling user data to json: %v", err)
	}

	fmt.Printf("User created successfully!\n\n")
	fmt.Println(string(userJsonFormat))

	return nil
}

func HandleReset(s *structure.State, cmd structure.Command) error {
	err := s.Database.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting users: %v", err)
	}
	fmt.Printf("Users successfully reset! Command -> %s", cmd.Name)
	return nil
}

func HandleUsers(s *structure.State, cmd structure.Command) error {
	users, err := s.Database.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting all users for list: %v", err)
	}

	for _, value := range users {
		if value.Name == s.Config.CurrentUserName {
			fmt.Printf("%s (current)\n", value.Name)
		} else {
			fmt.Printf("%s\n", value.Name)
		}
	}

	return nil
}

func HandleAgg(s *structure.State, cmd structure.Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	feedJsonFormat, err := json.MarshalIndent(*feed, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling rss feed to json: %v", err)
	}
	println(string(feedJsonFormat))

	return nil
}
