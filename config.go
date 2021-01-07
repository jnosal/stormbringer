package stormbringer

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	fMakeMaster = flag.Bool("master", false, "make this node master if unable to connect to the cluster ip provided.")
	fMasterIp   = flag.String("master-ip", "127.0.0.1:8001", "ip address of any node to connnect")
	fPort       = flag.Int("port", 8001, "ip address to run this node on. default is 8001.")
)

type Config struct {
	MakeMaster bool
	MasterIp   string
	Port       int
}

func ConfigFromFlags() Config {
	flag.Parse()
	return Config{
		MakeMaster: *fMakeMaster,
		MasterIp:   *fMasterIp,
		Port:       *fPort,
	}
}

func ConfigFromFile(fileName string) Config {
	var c Config
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("unable to read configuration", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatal("unable to decode configuration", err)
	}
	return c
}

func GetEnv(key, absentValue string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return absentValue
	}
	return v
}
