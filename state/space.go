package state

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"github.com/pkg/errors"
	"time"
)

type Space struct {
	Planets []*Planet
	Ships   []*Ship
	Width   uint32
	Height  uint32
}

func NewSpace() Space {
	rand.Seed(time.Now().UTC().UnixNano())
	space := Space{
		Width: 1400,
		Height: 500,
	}

	for i := uint32(0); i < 150; i++ {
		space.CreateNewPlanet();
	}

	return space
}

func (space *Space) CreateShip(planet *Planet) Ship {
	var id uint64 = 0
	if space.Ships != nil {
		id = space.Ships[len(space.Ships)-1].Id
		id++
	}

	s := Ship{
		Id: id,
	}
	space.Ships = append(space.Ships, &s)
	planet.orbiting = append(planet.orbiting, &s)
	return s
}

func (space *Space) CreateNewPlanet() *Planet {
	var id uint32 = 0
	if space.Planets != nil {
		id = space.Planets[len(space.Planets)-1].id
		id++
	}

	var x,y uint32
	for true{
		x, y = rand.Uint32() % space.Width, rand.Uint32() % space.Height
		if x < 50 || x > space.Width - 50{
			continue
		}

		if y < 50 || y > space.Height - 50{
			continue
		}

		for _, planet := range space.Planets{
			if x == planet.X || y == planet.Y{
				continue
			}
		}
		break
	}
	planet := Planet{
		id: id,
		X:  x,
		Y:  y,
		Control: rand.Float32(),
	}

	space.Planets = append(space.Planets, &planet)

	return &planet
}

func (space *Space) Serialize() ([]byte, error) {
	out := pb.Space{
		Width:space.Width,
		Height:space.Height,
	}
	for _, planet := range space.Planets {
		out.Planets = append(out.Planets, planet.Serialize())
	}
	return proto.Marshal(&out)
}

func Deserialize(data *[]byte) (*Space, error) {
	in := pb.Space{}
	err := proto.Unmarshal(*data, &in)
	if err != nil {
		return nil, errors.Wrap(err, "Could not deserialize.")
	}

	space := Space{
		Width: in.Width,
		Height: in.Height,
	}

	for _, planet := range in.Planets {
		p := Planet{
			id: planet.Id,
			X: planet.PosX,
			Y: planet.PosY,
			Empire: planet.Empire,
			Control: planet.Control,
		}

		space.Planets = append(space.Planets, &p)
		for _, _ = range planet.Orbiting {
			space.CreateShip(&p)
		}
	}

	return &space, nil
}
