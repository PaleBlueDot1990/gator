package middleware

import (
	"context"

	"github.com/PaleBlueDot1990/gator/internal/command"
	"github.com/PaleBlueDot1990/gator/internal/config"
	"github.com/PaleBlueDot1990/gator/internal/database"
)

func LoggedIn(handler func (state *config.State, cmd command.Command, user database.User) error) func(*config.State, command.Command) error {
	return func(state *config.State, cmd command.Command) error {
		user, err := state.DbQueries.GetUser(context.Background(), state.Cfg.CURRENT_USER_NAME)
		if err != nil {
			return err
		}

		return handler(state, cmd, user)
	}
}