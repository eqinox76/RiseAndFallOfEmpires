package client

import (
	"net"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	pb "github.com/eqinox76/RiseAndFallOfEmpires/proto"
	"fmt"
)

type Client struct {
	conn net.Conn
}

func (client *Client) Connect() error {
	conn, err := net.Dial("tcp", "localhost:9076")

	if err != nil {
		if client.conn != nil {
			client.conn.Close()
		}
		client.conn = nil
		return err
	}

	fmt.Println("Connected to server")
	client.conn = conn
	return nil
}


func (client *Client) SendCommand(cmd *pb.Command) error {
	data, err := proto.Marshal(cmd)
	if err != nil{
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

	_, err := client.conn.Read(header)
	if err != nil {
		client.conn.Close()
		client.conn = nil
		return nil, err
	}

	l := binary.LittleEndian.Uint32(header)

	msgbuffer := make([]byte, l)
	_, err = client.conn.Read(msgbuffer)
	if err != nil {
		return nil, err
	}

	space := pb.Space{}
	err = proto.Unmarshal(msgbuffer, &space)
	if err != nil {
		return nil, err
	}

	return &space, nil
}
