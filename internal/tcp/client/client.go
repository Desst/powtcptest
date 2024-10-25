package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	v1 "powtcptest/internal/protocol/v1"
	"powtcptest/internal/tcp"
	"time"
)

type Client struct {
	serverAddr  string
	readTimeout time.Duration

	challengeSolver PoWChallengeSolver
}

func NewClient(serverAddr string, readTimeout time.Duration, challengeSolver PoWChallengeSolver) *Client {
	return &Client{
		serverAddr:      serverAddr,
		readTimeout:     readTimeout,
		challengeSolver: challengeSolver,
	}
}

func (c *Client) RequestWordOfWisdom(ctx context.Context) (string, error) {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return "", fmt.Errorf("error connecting to server: %w", err)
	}

	log.Printf("Connected to %s", c.serverAddr)

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			if err := conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
				return "", fmt.Errorf("error setting read deadline: %w", err)
			}

			// Read incoming message
			var msg v1.Message
			if err := json.NewDecoder(conn).Decode(&msg); err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					continue
				}

				return "", fmt.Errorf("error reading message: %w", err)
			}

			switch msg.Type {
			case v1.MessageNewChallenge:
				newChallengeMsg, ok := msg.TypedMessage.(v1.NewChallengeMessage)
				if !ok {
					return "", fmt.Errorf("unable assert new challenge message type")
				}

				log.Printf("Challenge received: %s", newChallengeMsg.Challenge)

				start := time.Now()
				solution, err := c.challengeSolver.SolveChallenge(ctx, newChallengeMsg.Challenge, newChallengeMsg.Difficulty)
				elapsed := time.Since(start)
				if err != nil {
					return "", fmt.Errorf("error solving challenge: %w", err)
				}

				log.Printf("Challenge solution took %.2f", elapsed.Seconds())

				msg := v1.NewSolvedChallengeMessage(newChallengeMsg.Challenge, solution)
				if err := tcp.SendMessage(msg, conn); err != nil {
					return "", fmt.Errorf("error sending solved challenge msg: %w", err)
				}
			case v1.MessageWordOfWisdom:
				wordOfWisdomMsg, ok := msg.TypedMessage.(v1.WordOfWisdomMessage)
				if !ok {
					return "", fmt.Errorf("unable assert new word of wisdom message type")
				}

				return wordOfWisdomMsg.WordOfWisdom, nil
			default:
				return "", fmt.Errorf("received message of unexpected type %d", msg.Type)
			}
		}
	}
}
