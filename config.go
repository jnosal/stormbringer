package stormbringer

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	fMakeMaster = flag.Bool("master", false, "make this node master")
	fMakeWorker = flag.Bool("worker", false, "make this node worker")
	fMasterIp   = flag.String("master-ip", "127.0.0.1:8001", "ip address of any node to connnect")
	fHost       = flag.String("host", "0.0.0.0", "host to run this node on. default is 0.0.0.0")
	fPort       = flag.Int("port", 8001, "port number to run this node on. default is 8001.")
	fRps        = flag.Int("rps", 10, "requests per second")
)

type Config struct {
	MakeMaster bool
	MakeWorker bool
	MasterIp   string
	Port       int
	Host       string
	Rps        int
}

func (config *Config) IsStandalone() bool {
	return !config.MakeMaster && !config.MakeWorker
}

func (config *Config) IsMaster() bool {
	return config.MakeMaster
}

func (config *Config) IsWorker() bool {
	return config.MakeWorker
}

func ConfigFromFlags() Config {
	flag.Parse()
	return Config{
		MakeMaster: *fMakeMaster,
		MakeWorker: *fMakeWorker,
		MasterIp:   *fMasterIp,
		Port:       *fPort,
		Host:       *fHost,
		Rps:        *fRps,
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
