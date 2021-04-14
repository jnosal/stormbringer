package main

import (
	"stormbringer"
)

func main() {
	config := stormbringer.ConfigFromFlags()
	stormbringer.Run(config, new(stormbringer.DummyAttack))
}