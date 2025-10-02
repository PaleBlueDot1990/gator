package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/PaleBlueDot1990/gator/internal/command"
	"github.com/PaleBlueDot1990/gator/internal/config"
	"github.com/PaleBlueDot1990/gator/internal/database"
	"github.com/PaleBlueDot1990/gator/internal/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Creating Commands struct to hold a map of command <-> handler functions 
	cmds := command.Commands {
		HandlerMap: make(map[string]func(*config.State, command.Command) error),
	}
	cmds.Register("login", command.HandlerLogin)
	cmds.Register("register", command.HandlerRegister)
	cmds.Register("reset", command.HandleReset)
	cmds.Register("users", command.HandleUsers)
	cmds.Register("addfeed", middleware.LoggedIn(command.HandleAddFeed))
	cmds.Register("feeds", command.HandleFeeds)
	cmds.Register("follow", middleware.LoggedIn(command.HandleFollow))
	cmds.Register("following", middleware.LoggedIn(command.HandleFollowing))
	cmds.Register("agg", command.HandleAgg)

	// Getting the command from the command line interface 
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("program expects command name and possible arguments\n")
		os.Exit(1)
	}
	cmd := command.Command {
		Name: args[1],
		Args: args[2:],
	}

	// Getting configuration (username and db connection string)
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading configuration- %v\n", err)
		return
	}

	// Connecting to the database and getting all created db queries 
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Printf("error connecting to database- %v\n", err)
		return 
	}
	dbQueries := database.New(db)

	// Storing configuration and dbQueries in the current state 
	st := config.State {
		DbQueries: dbQueries,
		Cfg : cfg,
	}

	// Run the handler function for the given command 
	err = cmds.HandlerMap[cmd.Name](&st, cmd)
	if err != nil {
		fmt.Printf("unable to execute the command- %v\n", err)
		os.Exit(1)
	}
}