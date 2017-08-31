package main

import (
	"net"
	"encoding/binary"
	"net/http"
	"github.com/ajstarks/svgo"
	"fmt"
	"time"
	"math/rand"
)

var renderedWorld *[]byte

func main() {
	conn, err := net.Dial("tcp", "localhost:9076")
	if err != nil {
		panic(err)
	}

	go func() {
		http.HandleFunc("/space.svg", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	//go func() {
	//	for true {
	//		b := bytes.Buffer{}
	//
	//		writer := bufio.NewWriter(&b)
	//		dc := gg.NewContext(1000, 1000)
	//		dc.DrawCircle(500, 500, 200+rand.Float64()*200)
	//		dc.SetRGB(0, 0, 0)
	//		dc.Fill()
	//		dc.EncodePNG(writer)
	//		result := b.Bytes()
	//		renderedWorld = &result
	//	}
	//}()

	header := make([]byte, 4)
	for true {
		// read len

		_, err := conn.Read(header)
		if err != nil {
			panic(err)
		}
		l := binary.LittleEndian.Uint32(header)

		msgbuffer := make([]byte, l)
		_, err = conn.Read(msgbuffer)
		if err != nil {
			panic(err)
		}
	}
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

	writer.Header().Set("Content-Type", "image/svg+xml")                            // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache

	width := 500
	height := 500
	canvas := svg.New(writer)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 50 + rand.Int() % 100)
	canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()

	fmt.Println(time.Now().Sub(start))
}
