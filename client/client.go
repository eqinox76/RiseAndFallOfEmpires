package client

import (
	"io"
	"github.com/golang/protobuf/proto"
	"net"
	"fmt"
	"encoding/binary"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"time"
)

type Client struct {
	conn net.Conn
}

func (client *Client) Connect() error {
	conn, err := net.Dial("tcp", "localhost:9076")

	if err != nil {
		client.Close()
		return err
	}

	fmt.Println("Connected to server")
	client.conn = conn
	return nil
}

func (client *Client) Close() {
	if client.conn != nil {
		client.conn.Close()
		client.conn = nil
	}
}

func (client *Client) SendCommand(cmd *pb.Command) error {
	if client == nil || client.conn == nil {
		fmt.Errorf("connection already closed")
	}

	data, err := proto.Marshal(cmd)
	if err != nil {
		return err
	}

	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	_, err = client.conn.Write(append(length, data...))
	return err
}

func (client *Client) PollState() (*pb.Space, error) {
	// parse data and render it
	header := make([]byte, 4)

	// read len
	if client.conn == nil {
		err := client.Connect()
		if err != nil {
			return nil, err
		}
	}

	client.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	start := time.Now()
	_, err := client.conn.Read(header)
	io_time := time.Now().Sub(start)
	if err != nil {
		client.Close()
		return nil, err
	}

	l := binary.LittleEndian.Uint32(header)

	msgbuffer := make([]byte, l)
	start = time.Now()
	_, err = io.ReadFull(client.conn, msgbuffer)
	io_time += time.Now().Sub(start)
	if err != nil {
		return nil, err
	}

	space := pb.Space{}
	start = time.Now()
	err = proto.Unmarshal(msgbuffer, &space)
	unmarshal_time := time.Now().Sub(start)
	if err != nil {
		return nil, err
	}

	fmt.Printf("unmarshalled %d bytes in %s s with io %s s\n", l, unmarshal_time, io_time)

	return &space, nil
}
