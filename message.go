package stormbringer

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

/*
	Carries original byte message, optionally connection it came from,
	allows to send and receive maps along socket connection
*/
type TCPMessage struct {
	payload []byte
	conn    net.Conn
}

func (m *TCPMessage) getMap() (data map[string]interface{}) {
	err := json.Unmarshal(m.payload, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	return data
}

func NewTCPMessageFromMap(payload map[string]interface{}) (message *TCPMessage) {
	text, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Cannot send message: %s", payload)
		return
	}

	message = &TCPMessage{
		payload: text,
	}

	return
}

func NewTCPMessage(data []byte) (message *TCPMessage) {
	message = &TCPMessage{
		payload: data,
	}

	return
}
