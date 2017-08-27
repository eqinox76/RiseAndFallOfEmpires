package engine

import (
	"github.com/eqinox76/RiseAndFallOfEmpires/model"
	"math/rand"
)

func Step(space *model.Space) {
	for _, planet := range space.Planets {
		if rand.Float32() <= planet.Control {
			space.CreateShip(planet)
		}
	}
	return
}
