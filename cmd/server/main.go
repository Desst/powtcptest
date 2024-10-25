package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	serverCfg "powtcptest/internal/config/server"
	"powtcptest/internal/pow"
	"powtcptest/internal/quotes"
	"powtcptest/internal/tcp/server"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := serverCfg.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", err))
	}

	quotesCollection := quotes.NewCollection()
	powService := pow.NewService()

	tcpServer := server.NewServer(cfg.ListenAddr, time.Duration(cfg.SocketReadTimeoutSec)*time.Second,
		cfg.ChallengeDifficulty, quotesCollection, powService)

	if err := tcpServer.Start(); err != nil {
		log.Fatal(fmt.Errorf("unable to start tcp server: %w", err))
	}

	<-ctx.Done()
	toCtx, toCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer toCancel()

	if err := tcpServer.Shutdown(toCtx); err != nil {
		log.Fatal(fmt.Errorf("unable to shutdown tcp server: %w", err))
	}

	log.Println("Server stopped")
}
