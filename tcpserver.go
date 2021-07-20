package stormbringer

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"
)

const (
	CONNECTION_READ_WRITE_TIMEOUT = 60 * 5
	MAX_LINE_BYTES                = 1024
)

type TCPServer struct {
	address  string
	messages chan TCPMessage
	listener net.Listener
	chStop   chan interface{}
}

func resetReadDeadline(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(time.Second * CONNECTION_READ_WRITE_TIMEOUT))
	if err != nil {
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
	defer func() {
		_ = conn.Close()
		log.Printf("Connection from %v closed", conn.RemoteAddr())
	}()
	log.Printf("Got new TCP connection: %v", conn.RemoteAddr())
	done := make(chan struct{})

	// time out in CONNECTION_READ_WRITE_TIMEOUT seconds from now
	// if not data is received
	resetReadDeadline(conn)

	go func() {
		reader := &io.LimitedReader{
			R: conn,
			N: MAX_LINE_BYTES,
		}
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

			// reset number of bytes remaining in LimitReader
			reader.N = MAX_LINE_BYTES
			// reset read deadline
			resetReadDeadline(conn)
		}
		done <- struct{}{}
	}()
	<-done
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
			select {
			case <-server.chStop:
				return
			default:
				log.Print(err.Error())
				break
			}

			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *TCPServer) Stop() {
	log.Printf("Shutting down TCP server")
	close(server.chStop) // to tell server not to accept connections from listener
	server.listener.Close()
}

func NewTCPServer(address string) (server *TCPServer) {
	server = &TCPServer{
		address:  address,
		messages: make(chan TCPMessage),
		chStop:   make(chan interface{}),
	}

	return
}
