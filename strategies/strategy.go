package strategies

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/commands"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
)

type Strategy interface {
	Init(empire *state.Empire)

	Commands(space *state.Space) []commands.Command
}
