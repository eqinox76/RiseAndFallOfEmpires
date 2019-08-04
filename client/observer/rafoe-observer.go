package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/eqinox76/RiseAndFallOfEmpires/engine"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/ajstarks/svgo"
	"github.com/eqinox76/RiseAndFallOfEmpires/state"
	v "github.com/eqinox76/RiseAndFallOfEmpires/vector"
)

var picChan = make(chan []byte)

var freeColors = []string{
	"aliceblue",
	"antiquewhite",
	"aqua",
	"aquamarine",
	"azure",
	"beige",
	"bisque",
	"black",
	"blanchedalmond",
	"blue",
	"blueviolet",
	"brown",
	"burlywood",
	"cadetblue",
	"chartreuse",
	"chocolate",
	"coral",
	"cornflowerblue",
	"cornsilk",
	"crimson",
	"cyan",
	"darkblue",
	"darkcyan",
	"darkgoldenrod",
	"darkgray",
	"darkgreen",
	"darkkhaki",
	"darkmagenta",
	"darkolivegreen",
	"darkorange",
	"darkorchid",
	"darkred",
	"darksalmon",
	"darkseagreen",
	"darkslateblue",
	"darkslategray",
	"darkturquoise",
	"darkviolet",
	"deeppink",
	"deepskyblue",
	"dimgray",
	"dodgerblue",
	"firebrick",
	"floralwhite",
	"forestgreen",
	"fuchsia",
	"gainsboro",
	"ghostwhite",
	"gold",
	"goldenrod",
	"green",
	"greenyellow",
	"honeydew",
	"hotpink",
	"indianred",
	"indigo",
	"ivory",
	"khaki",
	"lavender",
	"lavenderblush",
	"lawngreen",
	"lemonchiffon",
	"lightblue",
	"lightcoral",
	"lightcyan",
	"lightgoldenrodyellow",
	"lightgray",
	"lightgreen",
	"lightpink",
	"lightsalmon",
	"lightseagreen",
	"lightskyblue",
	"lightslategray",
	"lightsteelblue",
	"lightyellow",
	"lime",
	"limegreen",
	"linen",
	"magenta",
	"maroon",
	"mediumaquamarine",
	"mediumblue",
	"mediumorchid",
	"mediumpurple",
	"mediumseagreen",
	"mediumslateblue",
	"mediumspringgreen",
	"mediumturquoise",
	"mediumvioletred",
	"midnightblue",
	"mintcream",
	"mistyrose",
	"moccasin",
	"navajowhite",
	"navy",
	"oldlace",
	"olive",
	"olivedrab",
	"orange",
	"orangered",
	"orchid",
	"palegoldenrod",
	"palegreen",
	"paleturquoise",
	"palevioletred",
	"papayawhip",
	"peachpuff",
	"peru",
	"pink",
	"plum",
	"powderblue",
	"purple",
	"red",
	"rosybrown",
	"royalblue",
	"saddlebrown",
	"salmon",
	"sandybrown",
	"seagreen",
	"seashell",
	"sienna",
	"silver",
	"skyblue",
	"slateblue",
	"slategray",
	"snow",
	"springgreen",
	"steelblue",
	"tan",
	"teal",
	"thistle",
	"tomato",
	"turquoise",
	"violet",
	"wheat",
	"white",
	"whitesmoke",
	"yellow",
	"yellowgreen}"}

type ObserverState struct {
	eng    engine.GameEngine
	colors map[*state.Empire]string
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()

	// start server
	go func() {
		http.HandleFunc("/space.svg", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	space := state.NewSpace(10)
	s := ObserverState{eng: engine.GameEngine{Space: &space}}
	s.eng.Init()

	s.colors = make(map[*state.Empire]string)
	for _, empire := range s.eng.Space.Empires {
		if empire.Passive{
			s.colors[empire] = "grey"
			continue
		}

		pos := rand.Intn(len(freeColors))
		s.colors[empire] = freeColors[pos]
		freeColors[pos] = freeColors[len(freeColors)-1 ]
		freeColors = freeColors[:len(freeColors)-1]
	}

	for ! s.eng.Space.Won() {

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		s.render(writer)

		picChan <- b.Bytes()

		s.eng.Step()

		time.Sleep(50 * time.Millisecond)
	}
}

func (obsState *ObserverState) render(writer *bufio.Writer) {
	width := int(obsState.eng.Space.Width)
	height := int(obsState.eng.Space.Height)

	canvas := svg.New(writer)
	canvas.Start(width, height)

	connected := make(map[*state.Planet]map[*state.Planet]bool)
	for _, planet := range obsState.eng.Space.Planets {
		// render connection
		connected_ids := connected[planet]
		for _, other := range planet.Connected {
			if connected_ids[other] {
				// already painted from the other side
				continue
			}

			// mark as painted
			_, ok := connected[other]
			if !ok {
				connected[other] = make(map[*state.Planet]bool)
			}
			connected[other][planet] = true
			empire := planet.Empire
			if planet.Empire == other.Empire && !empire.Passive {
				canvas.Line(int(planet.PosX), int(planet.PosY), int(other.PosX), int(other.PosY), fmt.Sprintf("stroke:white; stroke-width:2; stroke-opacity: 0.4; stroke: %s", obsState.colors[empire]))
			} else {
				canvas.Line(int(planet.PosX), int(planet.PosY), int(other.PosX), int(other.PosY), "stroke:white; stroke-width:2; stroke-opacity: 0.4")
			}
		}

		// render planet
		color := obsState.colors[planet.Empire]
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill-opacity: %f; fill: %s", planet.Control, color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 8, fmt.Sprintf("fill: none; stroke: %s; stroke-width: 2", color))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 10, fmt.Sprintf("fill: none; stroke: white; stroke-width: 1"))

		// #fleets
		//canvas.Text(int(planet.PosX), int(planet.PosY)-25, fmt.Sprint(len(planet.Fleets)), "text-anchor:middle;font-size:10px;fill:blue")

		fleets := planet.Fleets
		if len(fleets) == 1 {
			fleet := fleets[0]
			if !fleet.Empire.Passive {
				canvas.Text(int(planet.PosX), int(planet.PosY)-25, fmt.Sprint(fleet.Size()), fmt.Sprintf("text-anchor:middle;dominant-baseline;font-size:16px;stroke:%s;fill:%s", obsState.colors[fleet.Empire], obsState.colors[fleet.Empire]))
			}
		} else {
			center := v.Vec{
				X: float64(planet.PosX),
				Y: float64(planet.PosY),
			}
			nextDegree := 0.
			// more than one fleet. compute the position and rotation

			for _, fleet := range planet.Fleets {
				pos := center.MoveDegree(nextDegree, 25)

				canvas.Text(int(pos.X), int(pos.Y), fmt.Sprint(fleet.Size()), fmt.Sprintf("text-anchor:middle;dominant-baseline:central;font-size:16px;stroke:%s;fill:%s", obsState.colors[fleet.Empire], obsState.colors[fleet.Empire]))

				nextDegree += float64(360 / len(fleets))
			}
		}

		// show at most 10 ships per fleet
		for _, fleet := range planet.Fleets {
			size := int(math.Min(float64(fleet.Size()), 10))
			for i := 0; i <size; i++ {
				degree := 2 * math.Pi * rand.Float64()
				canvas.Circle(int(float64(planet.PosX)+(14*math.Sin(degree))),
					int(float64(planet.PosY)+(14*math.Cos(degree))),
					1,
					fmt.Sprintf("stroke: %s; stroke-width: 1", obsState.colors[fleet.Empire]))
			}
		}
	}

	canvas.Text(0, 10, fmt.Sprintf("Created: %s", time.Now()), "font-size:10px;fill:green")
	canvas.Text(0, 20, fmt.Sprintf("Round: %d", obsState.eng.Space.Round), "font-size:10px;fill:green")

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
	// firefox really needs this. no clue why the Cache-Control Headers are not enough.
	img.src = 'space.svg?t=t'+ Math.random(5);
	img.onload = function(){
		while (space.firstChild) {
			space.removeChild(space.firstChild);
		}

		space.insertBefore(img, null);
	}

	img.onerror = function() {
		console.log("Error loading " + this.src);
	};

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
	writer.Header().Set("Expires", "-1")

	//fmt.Println(len(picture), "bytes svg send")
	writer.Write(picture)
}
