package main

import (
	"net"
	"encoding/binary"
	"net/http"
	"github.com/fogleman/gg"
	"fmt"
	"bytes"
	"bufio"
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
		http.HandleFunc("/space.png", worldViewer)
		http.HandleFunc("/main", menuViewer)
		http.ListenAndServe(":8079", nil)
	}()

	go func() {
		for true {
			b := bytes.Buffer{}

			writer := bufio.NewWriter(&b)
			dc := gg.NewContext(1000, 1000)
			dc.DrawCircle(500, 500, rand.Float64()*500)
			dc.SetRGB(0, 0, 0)
			dc.Fill()
			dc.EncodePNG(writer)
			result := b.Bytes()
			renderedWorld = &result
		}
	}()

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

func menuViewer(writer http.ResponseWriter, request *http.Request){
data := `<html>

<head>

<script type="text/JavaScript">

var timeoutPeriod = 1000;
var imageURI = 'space.png';
var img = new Image();
img.onload = function() {
    var canvas = document.getElementById("x");
    var context = canvas.getContext("2d");

    context.drawImage(img, 0, 0);
    setTimeout(timedRefresh,timeoutPeriod);
};

function timedRefresh() {
    // just change src attribute, will always trigger the onload callback
    img.src = imageURI + '?d=' + Date.now();
}

</script>

<title>JavaScript Refresh Example</title>

</head>

<body onload="JavaScript:timedRefresh(1000);">

<canvas id="x" width="600" height="600" />

</body>
</html>`

	writer.Write([]byte(data))
}

func worldViewer(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()

	pic := *renderedWorld
	writer.Header().Set("Content-Type", "image/png") // set the content-type header
	writer.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // force no cache
	writer.Write(pic)
	fmt.Println(time.Now().Sub(start))
}
