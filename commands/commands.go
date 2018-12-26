package commands

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

//TODO make a interface with validate() and execute()

type Command interface {
	Execute()
	Validate() bool
}

type MoveCommand struct {
	Destination *state.Planet
	Fleet       *state.Fleet
}

func (moveCommand MoveCommand) Execute() {

	moveCommand.Fleet.Move(moveCommand.Destination)
}

func (moveCommand MoveCommand) Validate() bool {
	for _, pl := range moveCommand.Fleet.Position.Connected {
		if pl == moveCommand.Destination {
			return true
		}
	}
	return false
}
