package client

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerAddr           string
	SocketReadTimeoutSec int
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	config.ServerAddr = os.Getenv("SERVER_ADDR")

	readTimeoutSec, err := strconv.Atoi(os.Getenv("SOCKET_READ_TIMEOUT_SEC"))
	if err != nil {
		return nil, fmt.Errorf("could not parse SOCKET_READ_TIMEOUT_SEC: %w", err)
	}
	config.SocketReadTimeoutSec = readTimeoutSec
	return config, nil
}
