package main

import (
	"fmt"
	"stormbringer"
	"time"
)

func main() {
	client := stormbringer.NewTCPClient("localhost:8001")
	client.Connect()
	message := map[string]interface{}{"asd": 2, "bsd": "trhee"}
	times := []int{1, 2, 3, 4, 5, 6}
	for range times {
		err := client.Send(message)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 3)
	}
}
