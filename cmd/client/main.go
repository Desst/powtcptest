package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	clientCfg "powtcptest/internal/config/client"
	"powtcptest/internal/pow"
	tcpClient "powtcptest/internal/tcp/client"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := clientCfg.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", err))
	}

	powService := pow.NewService()

	client := tcpClient.NewClient(cfg.ServerAddr, time.Duration(cfg.SocketReadTimeoutSec)*time.Second, powService)

	quote, err := client.RequestWordOfWisdom(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to receive word of wisdom: %w", err))
	}

	log.Printf("Word of Wisdom received: \"%s\"", quote)
}
