package stormbringer

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
	if err != nil {
		log.Fatalf("Cannot connect to TCP server at: %s", client.address)
		return
	}

	log.Printf("Conected to TCP server: %s", client.address)

	client.connection = connection
}

func (client *TCPClient) Close() {
	_ = client.connection.Close()
}

func (client *TCPClient) Send(message map[string]interface{}) {
	tcpMessage := NewTCPMessageFromMap(message)

	_, err := client.connection.Write(tcpMessage.payload)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	_, err = client.connection.Write([]byte("\n"))
	if err != nil {
		log.Printf(err.Error())
		return
	}
}

func NewTCPClient(address string) (client *TCPClient) {
	client = &TCPClient{
		address: address,
	}

	return
}
