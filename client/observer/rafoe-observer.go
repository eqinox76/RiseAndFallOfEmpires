package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/ajstarks/svgo"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"github.com/gogo/protobuf/proto"
)

var path = flag.String("path", "state.tmp", "state which will be read")

var picChan = make(chan []byte)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	// start server
	go func() {
		http.HandleFunc("/space.svg", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	f, err := os.Open(*path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	reader, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}

	for {

		var l uint32
		err := binary.Read(reader, binary.LittleEndian, &l)
		if err != nil {
			fmt.Println("could not read length:", err)
			break
		}
		fmt.Println("going to read", l, "bytes of data")
		data := make([]byte, l, l)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			panic(err)
		}

		space := &pb.Space{}
		err = proto.Unmarshal(data, space)
		if err != nil {
			panic(err)
		}

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		render(writer, space)

		picChan <- b.Bytes()
	}
}

func render(writer *bufio.Writer, space *pb.Space) {
	width := int(space.Width)
	height := int(space.Height)

	canvas := svg.New(writer)
	canvas.Start(width, height)
	connected := make(map[uint32]map[uint32]bool)
	for _, planet := range space.Planets {
		// render connection
		connected_ids := connected[planet.Id]
		for _, other := range planet.Connected {
			if connected_ids[other] {
				// already painted from the other side
				continue
			}

			// mark as painted
			_, ok := connected[other]
			if !ok {
				connected[other] = make(map[uint32]bool)
			}
			connected[other][planet.Id] = true
			empire := space.Empires[planet.Empire]
			if planet.Empire == space.Planets[other].Empire && !empire.Passive {
				canvas.Line(int(planet.PosX), int(planet.PosY), int(space.Planets[other].PosX), int(space.Planets[other].PosY), fmt.Sprintf("stroke:white; stroke-width:2; stroke-opacity: 0.4; stroke: %s", empire.Color))
			} else {
				canvas.Line(int(planet.PosX), int(planet.PosY), int(space.Planets[other].PosX), int(space.Planets[other].PosY), "stroke:white; stroke-width:2; stroke-opacity: 0.4")
			}
		}

		// render planet
		color := space.Empires[planet.Empire].Color
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill-opacity: %f; fill: %s", planet.Control, color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill: none; stroke: %s; stroke-width: 4", color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 10, fmt.Sprintf("fill: none; stroke: white; stroke-width: 1"))

		// id
		//canvas.Text(int(planet.PosX), int(planet.PosY)-25, fmt.Sprint(planet.Id), "text-anchor:middle;font-size:10px;fill:blue")

		fleets := state.GetFleets(space.Ships, planet)
		if len(fleets) == 1 {
			for empire, fleet := range fleets {
				if !space.Empires[empire].Passive {
					canvas.Text(int(planet.PosX), int(planet.PosY)-25, fmt.Sprint(len(fleet)), fmt.Sprintf("text-anchor:middle;dominant-baseline;font-size:16px;stroke:%s;fill:%s", space.Empires[empire].Color, space.Empires[empire].Color))
				}
			}
		} else {
			center := v.Vec{
				X: float64(planet.PosX),
				Y: float64(planet.PosY),
			}
			nextDegree := 0.
			// more than one fleet. compute the position and rotation
			empires := sort.IntSlice{}
			for e := range fleets {
				empires = append(empires, int(e))
			}
			sort.Sort(empires)

			for _, eId := range empires {
				empire := space.Empires[uint32(eId)]
				fleet := fleets[uint32(eId)]
				pos := center.MoveDegree(nextDegree, 25)

				canvas.Text(int(pos.X), int(pos.Y), fmt.Sprint(len(fleet)), fmt.Sprintf("text-anchor:middle;dominant-baseline:central;font-size:16px;stroke:%s;fill:%s", empire.Color, empire.Color))

				nextDegree += float64(360 / len(fleets))
			}
		}

		// show at most 50 ships
		counter := 0
		for ship := range planet.Orbiting {
			if counter > 50 {
				break
			}
			counter++
			degree := 2 * math.Pi * rand.Float64()
			canvas.Circle(int(float64(planet.PosX)+(14*math.Sin(degree))),
				int(float64(planet.PosY)+(14*math.Cos(degree))),
				1,
				fmt.Sprintf("stroke: %s; stroke-width: 1", space.Empires[space.Ships[ship].Empire].Color))
		}
	}

	canvas.Text(0, 10, fmt.Sprintf("Created: %s", time.Now()), "font-size:10px;fill:green")
	canvas.Text(0, 20, fmt.Sprintf("Round: %d", space.Round), "font-size:10px;fill:green")

	canvas.End()
	writer.Flush()
}

func menuViewer(writer http.ResponseWriter, _ *http.Request) {
	data := `<html>

<head>

<title>JavaScript Refresh Example</title>

</head>

<body style="background-color:Black">

<div id="space" />

<script type="text/JavaScript">

var canvas = document.getElementById("space");
canvas.width  = window.innerWidth;
canvas.height = window.innerHeight;

function update() {

	var space = document.getElementById("space");
	var img = document.createElement('img');
	img.src = 'space.svg'
	img.onload = function(){
		while (space.firstChild) {
			space.removeChild(space.firstChild);
		}

		space.insertBefore(img, null);
	}

	// keep going
    requestAnimationFrame(update);
}

// schedule up the start
window.requestAnimationFrame(update);
</script>

</body>
</html>`

	writer.Write([]byte(data))
}

func worldViewer(writer http.ResponseWriter, _ *http.Request) {

	picture := <-picChan

	writer.Header().Set("Content-Type", "image/svg+xml")                        // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache

	fmt.Println(len(picture), "bytes svg send")
	writer.Write(picture)

}
