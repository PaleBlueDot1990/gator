package command

import (
	"context"
	"fmt"
	"time"

	"github.com/PaleBlueDot1990/gator/internal/config"
	"github.com/PaleBlueDot1990/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(state *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login command expects a single argument: the username")
	}

	username := cmd.Args[0]
	user, err := state.DbQueries.GetUser(context.Background(), username)
	if err != nil {
		return err 
	}

	err = state.Cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s have been logged in!\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Created Timestamp: %s\n", user.CreatedAt)
	fmt.Printf("User Updated Timestamp: %s\n", user.UpdatedAt)
	return nil
}

func HandlerRegister(state *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register command expetcs a single argument: the username")
	}

	dbQueryArgs := database.CreateUserParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	}

	user, err := state.DbQueries.CreateUser(context.Background(), dbQueryArgs)
	if err != nil {
		return err
	}

	err = state.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been registered!\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Created Timestamp: %s\n", user.CreatedAt)
	fmt.Printf("User Updated Timestamp: %s\n", user.UpdatedAt)
	return nil
}