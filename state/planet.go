package state

type Planet struct {
	PosX       uint32
	PosY       uint32
	Fleets     []*Fleet
	Control    float32
	Empire     *Empire
	Connected  []*Planet
	Production float32
}

func (planet *Planet) AddFleet(fleet *Fleet) {
	planet.Fleets = append(planet.Fleets, fleet)
}

func (planet Planet) EmpireFleet(empire *Empire) *Fleet {
	for _, fleet := range planet.Fleets {
		if fleet.Empire == empire {
			return fleet
		}
	}
	return nil
}

func (planet *Planet) RemoveFleet(fleet *Fleet) {
	for i, v := range (planet.Fleets) {
		if v == fleet {
			planet.Fleets[i] = planet.Fleets[len(planet.Fleets)-1]
			planet.Fleets = planet.Fleets[:len(planet.Fleets)-1]
			return
		}
	}
	panic("Fleet foes not exist")
}
