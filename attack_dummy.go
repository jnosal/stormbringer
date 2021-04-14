package stormbringer

import (
    "log"
    "time"
)

type DummyAttack struct{}

func (a *DummyAttack) Setup(c Config) error {
	return nil
}

func (a *DummyAttack) Do() {
	time.Sleep(100 * time.Millisecond)
	log.Println("ZOMBIE!!!")
}

func (a *DummyAttack) Teardown() error {
	return nil
}

func (a *DummyAttack) BeforeRun(c Config) error {
	log.Println("before run")
	return nil
}


func (a *DummyAttack) AfterRun(c Config) error {
	log.Println("after run")
	return nil
}
