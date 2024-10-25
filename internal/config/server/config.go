package server

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ListenAddr           string
	ChallengeDifficulty  int
	SocketReadTimeoutSec int
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	config.ListenAddr = os.Getenv("LISTEN_ADDR")
	challengeDifficulty, err := strconv.Atoi(os.Getenv("CHALLENGE_DIFFICULTY"))
	if err != nil {
		return nil, fmt.Errorf("could not parse CHALLENGE_DIFFICULTY: %w", err)
	}
	config.ChallengeDifficulty = challengeDifficulty

	readTimeoutSec, err := strconv.Atoi(os.Getenv("SOCKET_READ_TIMEOUT_SEC"))
	if err != nil {
		return nil, fmt.Errorf("could not parse SOCKET_READ_TIMEOUT_SEC: %w", err)
	}
	config.SocketReadTimeoutSec = readTimeoutSec

	return config, nil
}
