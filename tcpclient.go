package stormbringer

import (
	"bufio"
	"errors"
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
	writer     *bufio.Writer
	connected  bool
}

func (client *TCPClient) Connect() {
	connection, err := net.DialTimeout("tcp", client.address, time.Duration(CONNECTION_TIMEOUT)*time.Second)
	if err != nil {
		log.Fatalf("Cannot connect to TCP server at: %s", client.address)
		return
	}

	log.Printf("Conected to TCP server: %s", client.address)

	// TODO: handle gently
	//err = connection.(*net.TCPConn).SetKeepAlive(true)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//err = connection.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 3)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	client.connected = true
	client.connection = connection
	client.writer = bufio.NewWriter(connection)

}

func (client *TCPClient) Close() {
	_ = client.connection.Close()
}

func (client *TCPClient) Send(message map[string]interface{}) error {
	// TODO: possibly lock/umlock
	if !client.connected {
		return errors.New("client is not connected")
	}
	tcpMessage := NewTCPMessageFromMap(message)

	if _, err := client.writer.Write(tcpMessage.payload); err != nil {
		return err
	}

	if _, err := client.writer.WriteString("\r\n"); err != nil {
		return err
	}

	if err := client.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func NewTCPClient(address string) (client *TCPClient) {
	client = &TCPClient{
		address: address,
	}

	return
}
