package stormbringer

import (
	"flag"
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
