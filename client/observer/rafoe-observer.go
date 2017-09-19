package main

import (
	"net/http"
	"github.com/ajstarks/svgo"
	"fmt"
	"time"
	"bytes"
	"bufio"
	"sync"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"math/rand"
	"math"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	"github.com/eqinox76/RiseAndFallOfEmpires/client"
	"sort"
	"strings"
)

var registered []chan []byte
var mux sync.Mutex

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// start server
	go func() {
		http.HandleFunc("/space.svg", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	// connect to the server and render the gamestate forever
	c := client.Client{}

	for {
		space, err := c.PollState()
		if err != nil{
			fmt.Errorf("%s while reading new state", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		start := time.Now()
		render(writer, space)
		fmt.Println(time.Now().Sub(start))

		mux.Lock()
		for _, c := range registered {
			data := b.Bytes()
			c <- data
			close(c)
		}
		registered = nil
		mux.Unlock()
	}
}

func render(writer *bufio.Writer, space *pb.Space) {
	width := int(space.Width)
	height := int(space.Height)

	canvas := svg.New(writer)
	canvas.Start(width, height)
	connected := make(map[uint32] map[uint32] bool)
	for _, planet := range space.Planets{
		// render connection
		connected_ids := connected[planet.Id]
		for _, other := range planet.Connected {
			if connected_ids[other]{
				// already painted from the other side
				continue
			}

			// mark as painted
			_, ok := connected[other]
			if !ok{
				connected[other] = make(map[uint32] bool)
			}
			connected[other][planet.Id] = true
			canvas.Line(int(planet.PosX), int(planet.PosY), int(space.Planets[other].PosX), int(space.Planets[other].PosY), "stroke:blue; stroke-width:2; stroke-opacity: 0.2")
		}

		// render planet
		color := space.Empires[planet.Empire].Color
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill-opacity: %f; fill: %s", planet.Control, color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill: none; stroke: %s; stroke-width: 4", color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 10, fmt.Sprintf("fill: none; stroke: black; stroke-width: 1"))

		canvas.Text(int(planet.PosX), int(planet.PosY)-25, fmt.Sprint(planet.Id), "text-anchor:middle;font-size:10px;fill:blue")

		fleets := state.GetFleets(space.Ships, planet)
		if planet.Empire != 0 {
			// show control
			// canvas.Text(int(planet.PosX), int(planet.PosY)+20, fmt.Sprint("Control:", planet.Control), "text-anchor:middle;font-size:10px;fill:green")
			// show ships
			canvas.Text(int(planet.PosX), int(planet.PosY)-15, fmt.Sprint(len(fleets[planet.Empire])), "text-anchor:middle;font-size:12px;stroke:black;stroke-width:0.25;fill:" + color)
		}


		if len(fleets) > 1{
			var text []string
			for key, value := range fleets {
				text = append(text, fmt.Sprintf("%s: %d,", space.Empires[key].Color, len(value)))
			}
			sort.Strings(text)
			canvas.Text(int(planet.PosX), int(planet.PosY)-25, strings.Join(text, ";"), "text-anchor:middle;font-size:10px;fill:green")
		}

		// show at most 50 ships
		counter := 0
		for ship, _ := range planet.Orbiting{
			if counter > 50{
				break
			}
			counter++
			degree := 2 * math.Pi * rand.Float64()
			canvas.Circle(int(float64(planet.PosX) + (14 * math.Sin(degree))),
				int(float64(planet.PosY) + (14 * math.Cos(degree))),
				1,
				fmt.Sprintf("stroke: %s; stroke-width: 1", space.Empires[space.Ships[ship].Empire].Color))
		}
	}


	canvas.Text(0, 10, fmt.Sprintf("Created: %s", time.Now()), "font-size:10px")

	canvas.End()
	writer.Flush()
}

func menuViewer(writer http.ResponseWriter, request *http.Request) {
	data := `<html>

<head>

<title>JavaScript Refresh Example</title>

</head>

<body>

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

func worldViewer(writer http.ResponseWriter, request *http.Request) {
	c := make(chan []byte)
	mux.Lock()
	registered = append(registered, c)
	mux.Unlock()

	writer.Header().Set("Content-Type", "image/svg+xml")                        // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache

	data := <- c
	writer.Write(data)

}
