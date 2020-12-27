package tcp

import (
	"log"
	"net"
	"time"
)

const (
	CONNECTION_TIMEOUT = 10
)

type TCPClient struct {
	address    string
	connection net.Conn
}

func (client *TCPClient) Connect() {
	connection, err := net.DialTimeout("tcp", client.address, time.Duration(CONNECTION_TIMEOUT)*time.Second)
	client.connection = connection

	if err != nil {
		log.Printf("Cannot connect to TCP server at: %s", client.address)
		return
	}
}

func (client *TCPClient) Close() {
	_ = client.connection.Close()
}

func (client *TCPClient) Send(message map[string]interface{}) {
	tcpMessage := NewTCPMessageFromMap(message)

	client.connection.Write(tcpMessage.payload)
	client.connection.Write([]byte("\n"))
}

func NewTCPClient(address string) (client *TCPClient) {
	client = &TCPClient{
		address: address,
	}

	return
}
