package command

import (
	"github.com/PaleBlueDot1990/gator/internal/config"
)

type Command struct {
	Name string
	Args []string 
}

type Commands struct {
	HandlerMap map[string]func(*config.State, Command) error 
}

func (cmds* Commands) Register(name string, f func(*config.State, Command) error) {
	cmds.HandlerMap[name] = f
}

func (cmds* Commands) Run(s* config.State, cmd Command) error {
	f := cmds.HandlerMap[cmd.Name]
	return f(s, cmd)
}