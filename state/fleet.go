package state

type Fleet struct {
	Empire       *Empire
	Position     *Planet
	LightSquads  int
	HeavySquads  int
	RangedSquads int
}

func (fleet Fleet) Size() int {
	return fleet.LightSquads + fleet.HeavySquads + fleet.RangedSquads
}

func (fleet *Fleet) Move(planet *Planet) {
	fleet.Position.RemoveFleet(fleet)
	fleet.Position = planet
	planet.AddFleet(fleet)
}
