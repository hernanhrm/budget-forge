package localconfig

import (
	"os"
	"strconv"
)

type LocalConfig struct {
	Service  Service
	Database Database
}

type Service struct {
	Port  func() int
	Name  string
	Debug bool
}

type Database struct {
	URL string
}

func getPort() int {
	p := os.Getenv("SERVICE_PORT")
	if p == "" {
		p = os.Getenv("PORT")
	}
	if p == "" {
		return 8080
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 8080
	}
	return port
}

func getDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
