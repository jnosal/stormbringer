package stormbringer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	STATE_INITIAL = "INTITIAL"
)

type Engine struct {
	config Config
	state  string
	wg     sync.WaitGroup
	chDone chan struct{}
}

func (engine *Engine) startTCPServer() {
	server := NewTCPServer(fmt.Sprintf("%s:%d", engine.config.Host, engine.config.Port))
	server.Start()
}

func (engine *Engine) startTCPClient() {
	client := NewTCPClient(engine.config.MasterIp)
	client.Connect()
}

func Run(config Config) {
	log.Printf("Starting engine using: %+v", config)

	engine := Engine{
		state:  STATE_INITIAL,
		config: config,
		chDone: make(chan struct{}),
	}

	if engine.config.IsMaster() {
		go engine.startTCPServer()
	}

	if engine.config.IsWorker() {
		go engine.startTCPClient()
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received %s. Stopping", sig)
		return
	}

}
