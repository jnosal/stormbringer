package stormbringer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
	config Config
	state  string
	wg     sync.WaitGroup
	chDone chan bool
	chNext chan bool
}

func (engine *Engine) startMaster() {
	server := NewTCPServer(fmt.Sprintf("%s:%d", engine.config.Host, engine.config.Port))
	server.Start()
}

func (engine *Engine) startWorker() {
	client := NewTCPClient(engine.config.MasterIp)
	client.Connect()
}

func (engine *Engine) scheduleAttacks(config Config, attack Attack) {
	if v, ok := attack.(BeforeRunner); ok {
		if err := v.BeforeRun(config); err != nil {
			log.Printf("BeforeRun failed:%v\n", err)
		}
	}

	limiter := time.Tick(time.Millisecond * 100)

	for {
		select {
		case <-engine.chDone:
			log.Print("Stopping...")
			return
		default:
			<-limiter
			go attack.Do()
		}
	}
	if v, ok := attack.(AfterRunner); ok {
		if err := v.AfterRun(config); err != nil {
			log.Printf("AfterRun failed:%v\n", err)
		}
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
		engine.scheduleAttacks(config, attack)
	}

	if engine.config.IsMaster() {
		go engine.startMaster()
	}

	if engine.config.IsWorker() {
		go engine.startWorker()
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received %s. Stopping", sig)
		return
	}

}
