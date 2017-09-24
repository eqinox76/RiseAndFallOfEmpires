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
	"google.golang.org/grpc"
	"log"
	"golang.org/x/net/context"
	"io"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
	"sort"
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
	conn, err := grpc.Dial("localhost:9076", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewGameServerClient(conn)

	stream, err := c.Observe(context.Background(), &pb.ID{})
	if err != nil {
		panic(err)
	}

	for {
		space, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server closed connection")
			break
		} else if err != nil {
			log.Printf("%s while reading new state\n", err)
			break
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
			canvas.Line(int(planet.PosX), int(planet.PosY), int(space.Planets[other].PosX), int(space.Planets[other].PosY), "stroke:white; stroke-width:2; stroke-opacity: 0.4")
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
				if ! space.Empires[empire].Passive {
					canvas.Text(int(planet.PosX), int(planet.PosY)-20, fmt.Sprint(len(fleet)), "text-anchor:middle;dominant-baseline;font-size:12px;stroke:white;stroke-width:0.5;fill:"+space.Empires[empire].Color)
				}
			}
		} else {
			center := v.Vec{float64(planet.PosX), float64(planet.PosY)}
			nextDegree := 0.
			// more than one fleet. compute the position and rotation
			empires := sort.IntSlice{}
			for e, _ := range fleets {
				empires = append(empires, int(e))
			}
			sort.Sort(empires)

			fmt.Print(empires, fleets)
			for _, eId := range empires {
				empire := space.Empires[uint32(eId)]
				fleet := fleets[uint32(eId)]
				pos := center.MoveDegree(nextDegree, 20)

				canvas.Text(int(pos.X), int(pos.Y), fmt.Sprint(len(fleet)), fmt.Sprintf("text-anchor:middle;dominant-baseline:central;font-size:12px;stroke:white;stroke-width:0.5;fill:%s", empire.Color))

				nextDegree += float64(360 / len(fleets))
			}
		}

		// show at most 50 ships
		counter := 0
		for ship, _ := range planet.Orbiting {
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

func menuViewer(writer http.ResponseWriter, request *http.Request) {
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

func worldViewer(writer http.ResponseWriter, request *http.Request) {
	c := make(chan []byte)
	mux.Lock()
	registered = append(registered, c)
	mux.Unlock()

	writer.Header().Set("Content-Type", "image/svg+xml")                        // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache

	data := <-c
	writer.Write(data)

}
