package main

import (
	"net"
	"encoding/binary"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9076")
	if err != nil {
		panic(err)
	}

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
