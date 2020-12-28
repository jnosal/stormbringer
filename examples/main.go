package main

import (
	"fmt"
	"stormbringer"
	"strconv"
)

func main() {
	config := stormbringer.ConfigFromFlags()

	if config.MakeMaster {
		server := stormbringer.NewTCPServer(fmt.Sprintf(":%d", config.Port))
		server.Start()
	} else {
		client := stormbringer.NewTCPClient(config.MasterIp)
		client.Connect()
		for i := 0; i < 500; i++ {
			data := map[string]interface{}{
				"dupa": strconv.Itoa(i + 1),
				"el":   strconv.Itoa(i * 2),
			}
			client.Send(data)
		}
		client.Close()
	}
}
