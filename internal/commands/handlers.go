package commands

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AhmettCelik/blog-aggregator/internal/database"
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
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error: not enough arguments")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("Error parsing duration string to time.Duration value: %v", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", cmd.Args[1])

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
		fmt.Printf("Post created! Continuing..\n")
	}
}

func HandleAddFeed(s *structure.State, cmd structure.Command, user database.User) error {
	if len(cmd.Args) < 3 {
		return fmt.Errorf("Error: not enough arguments")
	}

	now := time.Now()
	feedName := cmd.Args[1]
	feedURL := cmd.Args[2]

	feedParams := database.CreateFeedParams{
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.Database.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error creating a feed: %v", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		CreatedAt: now,
		UpdatedAt: now,
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	_, err = s.Database.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		log.Fatalf("Error creating feed follow: %v", err)
		os.Exit(1)
	}

	feedJsonFormat, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling feed table struct to json: %v", err)
	}
	println(string(feedJsonFormat))

	return nil
}

func HandleFeeds(s *structure.State, cmd structure.Command) error {
	feedsInfos, err := s.Database.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds: %v", err)
	}

	for _, feedInfo := range feedsInfos {
		fmt.Printf("%s\n", feedInfo.Name)
		fmt.Printf("%s\n", feedInfo.Url)
		fmt.Printf("%s\n", feedInfo.UserName)
		fmt.Printf("\n")
	}

	return nil
}

func HandleFollow(s *structure.State, cmd structure.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error: you must provide a URL for this command")
	}

	url := cmd.Args[1]
	now := time.Now()
	feed, err := s.Database.GetFeedForUrl(context.Background(), url)
	if err != nil {
		log.Fatalf("Error getting feed for url: %v", err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		CreatedAt: now,
		UpdatedAt: now,
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	feedFollow, err := s.Database.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		log.Fatalf("Error creating feed follow: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Feed name: %s\n", feedFollow.FeedName)
	fmt.Printf("Current user name: %s\n", feedFollow.UserName)

	return nil
}

func HandleFollowing(s *structure.State, cmd structure.Command, user database.User) error {
	feedFollows, err := s.Database.GetFeetFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("Error getting feed follows: %v", err)
		os.Exit(1)
	}

	for _, f_f := range feedFollows {
		fmt.Printf("%s\n", f_f.FeedName)
	}

	return nil
}

func HandleUnfollow(s *structure.State, cmd structure.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		log.Fatal("You must provide a URL for use this command")
		os.Exit(1)
	}

	url := cmd.Args[1]
	feed, err := s.Database.GetFeedForUrl(context.Background(), url)
	if err != nil {
		log.Fatalf("Error getting feed for url: %v", err)
		os.Exit(1)
	}

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	s.Database.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)

	fmt.Printf("Unfollowing: %s\n", feed.Name)

	return nil
}

func HandleBrowse(s *structure.State, cmd structure.Command) error {
	var limit int32

	if len(cmd.Args) < 2 || cmd.Args[1] == "" {
		limit = 2
	} else {
		limit64, err := strconv.ParseInt(cmd.Args[1], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(limit64)
	}

	posts, err := s.Database.GetPostsForUser(context.Background(), limit)
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %v", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s", post.Title)
		fmt.Printf("Description: %s", post.Description)
	}

	return nil
}
