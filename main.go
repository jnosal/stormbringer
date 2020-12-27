package main

import (
	"flag"
	"fmt"
	"math/rand"
	"stormbringer/tcp"
	"strconv"
	"time"
)

func main() {
	makeMaster := flag.Bool("master", false, "make this node master if unable to connect to the cluster ip provided.")
	masterIp := flag.String("master-ip", "127.0.0.1:8001", "ip address of any node to connnect")
	port := flag.Int("port", 8001, "ip address to run this node on. default is 8001.")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	id := rand.Intn(99999999)

	fmt.Printf("Config: makeMaster=%t clusterIp=%s port=%d id=%d", *makeMaster, *masterIp, *port, id)
	fmt.Println()

	if *makeMaster {
		server := tcp.NewTCPServer(fmt.Sprintf(":%d", *port))
		server.Start()
	} else {
		client := tcp.NewTCPClient(*masterIp)
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
