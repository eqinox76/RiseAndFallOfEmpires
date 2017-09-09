package main

import (
	"net"
	"encoding/binary"
	"net/http"
	"github.com/ajstarks/svgo"
	"fmt"
	"time"
	"bytes"
	"bufio"
	"sync"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"github.com/golang/protobuf/proto"
)

var registered []chan []byte
var mux sync.Mutex

func main() {
	conn, err := net.Dial("tcp", "localhost:9076")
	if err != nil {
		panic(err)
	}

	// start server
	go func() {
		http.HandleFunc("/space.svg", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	// parse data and render it
	header := make([]byte, 4)
	for {
		// read len
		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:9076")
			if err != nil {
				fmt.Println(err)
				conn = nil
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}

		_, err := conn.Read(header)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
			conn.Close()
			conn = nil

			continue
		}
		l := binary.LittleEndian.Uint32(header)

		msgbuffer := make([]byte, l)
		_, err = conn.Read(msgbuffer)
		if err != nil {
			panic(err)
		}

		space := pb.Space{}
		err = proto.Unmarshal(msgbuffer, &space)
		if err != nil {
			panic(err)
		}

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		render(writer, &space)

		mux.Lock()
		for _, c := range (registered) {
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
	for _, planet := range space.Planets{
		// render connection
		for _, id := range planet.Connected {
			canvas.Line(int(planet.PosX), int(planet.PosY), int(space.Planets[id].PosX), int(space.Planets[id].PosY), "stroke:blue; stroke-width:2; stroke-opacity: 0.2")
		}

		// render planet
		color := space.Empires[planet.Empire].Color
		canvas.Circle(int(planet.PosX), int(planet.PosY), 10, fmt.Sprintf("fill: none; stroke: black; stroke-width: 1"))
		canvas.Circle(int(planet.PosX), int(planet.PosY), 10, fmt.Sprintf("fill-opacity: %f; fill: %s", planet.Control, color))
		if (planet.Empire != 0) {
			// show id
			canvas.Text(int(planet.PosX), int(planet.PosY)-20, fmt.Sprint("Id:", planet.Id), "text-anchor:middle;font-size:10px;fill:green")
			// show control
			canvas.Text(int(planet.PosX), int(planet.PosY)+20, fmt.Sprint("Control:", planet.Control), "text-anchor:middle;font-size:10px;fill:green")
		}
		if len(planet.Orbiting) > 0 {
			// show ships
			canvas.Text(int(planet.PosX), int(planet.PosY)-10, fmt.Sprint("#Ships:", len(planet.Orbiting)), "text-anchor:middle;font-size:10px;fill:green")
		}

	}

	canvas.End()
	writer.Flush()
}

func menuViewer(writer http.ResponseWriter, request *http.Request) {
	data := `<html>

<head>

<title>JavaScript Refresh Example</title>

</head>

<body>

<canvas id="space" />

<script type="text/JavaScript">

var canvas = document.getElementById("space");
canvas.width  = window.innerWidth;
canvas.height = window.innerHeight;

function update() {

	var c = document.getElementById("space");
	var ctx = c.getContext("2d");
	var img = new Image();
	img.src = 'space.svg'
	img.onload = function(){
		ctx.clearRect(0, 0, c.width, c.height);
		ctx.drawImage(img, 0, 0);
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
	start := time.Now()
	c := make(chan []byte)
	mux.Lock()
	registered = append(registered, c)
	mux.Unlock()

	writer.Header().Set("Content-Type", "image/svg+xml")                        // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache

	data := <- c
	writer.Write(data)

	fmt.Println(time.Now().Sub(start))
}
