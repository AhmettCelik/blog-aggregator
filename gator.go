package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AhmettCelik/blog-aggregator/internal/commands"
	"github.com/AhmettCelik/blog-aggregator/internal/config"
	"github.com/AhmettCelik/blog-aggregator/internal/database"
	"github.com/AhmettCelik/blog-aggregator/internal/middleware"
	"github.com/AhmettCelik/blog-aggregator/internal/structure"
	_ "github.com/lib/pq"
)

func startGator() {
	var cfg config.Config = config.Read()
	var state *structure.State = new(structure.State)
	var cmds *structure.Commands = new(structure.Commands)

	dbURL := cfg.DatabaseUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error: failed to connect to database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state.Database = dbQueries

	argsWithoutPath := os.Args[1:]
	if len(argsWithoutPath) == 0 {
		log.Fatal("error: not enought arguments were provided.")
	}

	state.Config = &cfg

	cmds.Handlers = make(map[string]func(*structure.State, structure.Command) error)

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandleReset)
	cmds.Register("users", commands.HandleUsers)
	cmds.Register("agg", commands.HandleAgg)
	cmds.Register("feeds", commands.HandleFeeds)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(commands.HandleAddFeed))
	cmds.Register("follow", middleware.MiddlewareLoggedIn(commands.HandleFollow))
	cmds.Register("following", middleware.MiddlewareLoggedIn(commands.HandleFollowing))
	cmds.Register("unfollow", middleware.MiddlewareLoggedIn(commands.HandleUnfollow))

	err = cmds.Run(state, structure.Command{Name: argsWithoutPath[0], Args: argsWithoutPath})
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
}
