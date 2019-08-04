package commands

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

type Command interface {
	Execute(space *state.Space)
	Validate() bool
}

type MoveCommand struct {
	Destination *state.Planet
	Fleet       *state.Fleet
}

func (moveCommand MoveCommand) Execute( *state.Space) {

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

type FleetMergeCommand struct {
	Fleet []*state.Fleet
}

func (fleetMergeCommand FleetMergeCommand) Execute(space *state.Space) {
	fst, targets := fleetMergeCommand.Fleet[0], fleetMergeCommand.Fleet[1:]
	for _, trg := range targets {
		fst.MergeFrom(trg)
		space.DestroyFleet(trg)
	}
}

func (fleetMergeCommand FleetMergeCommand) Validate() bool {
	var e *state.Empire
	for _, f := range fleetMergeCommand.Fleet {
		if e == nil || e == f.Empire {
			e = f.Empire
		} else {
			return false
		}
	}
	return true
}
