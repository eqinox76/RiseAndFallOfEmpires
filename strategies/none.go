package strategies

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

type NoneStrategy struct {
}

func (strategy NoneStrategy) Init(empire *state.Empire) {
}

func (strategy NoneStrategy) Commands(space *state.Space) []commands.Command {
	return nil
}
