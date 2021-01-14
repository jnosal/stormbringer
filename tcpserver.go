package stormbringer

import (
	"bufio"
	"log"
	"net"
	"time"
)

const (
	CONNECTION_READLINE_DEADLINE = 60 * 5
)

type TCPServer struct {
	address  string
	messages chan TCPMessage
	listener net.Listener
}

func increaseDeadline(conn net.Conn) {
	if err := conn.SetReadDeadline(time.Now().Add(time.Second * CONNECTION_READLINE_DEADLINE)); err != nil {
		log.Println("Cannot increase read deadline")
	}
}

func (server *TCPServer) handleMessages() {
	for {
		tcpMessage := <-server.messages
		log.Printf("Got string message: %s", tcpMessage.payload)
		log.Printf("Got map object: %s", tcpMessage.getMap())
	}
}

func (server *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Got new TCP connection: %v", conn.RemoteAddr())

	increaseDeadline(conn)
	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)

	for {
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
			}
			break
		}
		server.messages <- TCPMessage{scanner.Bytes(), conn}
		increaseDeadline(conn)
	}

	log.Printf("Connection from %v closed", conn.RemoteAddr())
}

func (server *TCPServer) Start() {
	_, err := net.ResolveTCPAddr("tcp4", server.address)
	if err != nil {
		log.Printf("Cannot start TCP server at: %s", server.address)
		return
	}

	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Printf("Cannot start TCP server at: %s", server.address)
		return
	}

	server.listener = listener

	defer listener.Close()

	go server.handleMessages()

	log.Printf("Starting TCP server at %s", server.address)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err.Error())
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *TCPServer) Stop() {
	log.Printf("Shutting down TCP server")
	server.listener.Close()
}

func NewTCPServer(address string) (server *TCPServer) {
	server = &TCPServer{
		address:  address,
		messages: make(chan TCPMessage),
	}

	return
}
