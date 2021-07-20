package stormbringer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	STATE_INITIAL = "INTITIAL"
)

type BeforeRunner interface {
	BeforeRun(c Config) error
}

type AfterRunner interface {
	AfterRun(c Config) error
}

type Engine struct {
	config    Config
	state     string
	tcpServer *TCPServer
	tcpClient *TCPClient
	chDone    chan bool
	chNext    chan bool
}

func (engine *Engine) startMaster() {
	server := NewTCPServer(fmt.Sprintf("%s:%d", engine.config.Host, engine.config.Port))
	engine.tcpServer = server
	server.Start()
}

func (engine *Engine) startWorker() {
	client := NewTCPClient(engine.config.MasterIp)
	engine.tcpClient = client
	client.Connect()
}

func (engine *Engine) scheduleAttacks(attack Attack) {
	if v, ok := attack.(BeforeRunner); ok {
		if err := v.BeforeRun(engine.config); err != nil {
			log.Printf("BeforeRun failed:%v\n", err)
		}
	}

	limiter := time.Tick(time.Millisecond * 100)

	for {
		select {
		default:
			<-limiter
			go attack.Do()
		}
	}
}

func (engine *Engine) Stop(attack Attack) {
	if v, ok := attack.(AfterRunner); ok {
		if err := v.AfterRun(engine.config); err != nil {
			log.Printf("AfterRun failed:%v\n", err)
		}
	}

	if engine.config.IsMaster() && engine.tcpServer != nil {
		engine.tcpServer.Stop()
	}

	if engine.config.IsWorker() && engine.tcpClient != nil {
		engine.tcpClient.Close()
	}
}

func Run(config Config, attack Attack) {
	log.Printf("Starting engine using: %+v", config)

	engine := Engine{
		state:  STATE_INITIAL,
		config: config,
		chDone: make(chan bool),
		chNext: make(chan bool),
	}

	if engine.config.IsStandalone() {
		go engine.scheduleAttacks(attack)
	}

	if engine.config.IsMaster() {
		go engine.startMaster()
	}

	if engine.config.IsWorker() {
		go engine.startWorker()
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigChan:
			log.Printf("Received %s. Stopping", sig)
			engine.Stop(attack)
			engine.chDone <- true
			return
		}
	}()
	<-engine.chDone

}
