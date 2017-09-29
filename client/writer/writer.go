package main

import (
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"google.golang.org/grpc"
	"log"
	"golang.org/x/net/context"
	"io"
	"flag"
	"os"
	"compress/gzip"
	"fmt"
)

var path = flag.String("write", ".", "were to write the data to")

var step = flag.Int("step", 100, "only every nth step will be persisted")

func main() {
	flag.Parse()

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

	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}

	writer, err := gzip.NewWriterLevel(f, gzip.DefaultCompression)
	if err != nil {
		panic(err)
	}

	defer writer.Close()
	defer f.Close()

	for {
		counter := 0
		space, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server closed connection")
			break
		} else if err != nil {
			log.Printf("%s while reading new state\n", err)
			break
		}

		if counter % *step == 0 {
			data, err := space.Marshal()
			writer.Write(data)
			fmt.Println(len(data), "bytes before compression")
			if err != nil {
				panic(err)
			}
		}
		counter++
	}
}
