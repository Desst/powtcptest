package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	v1 "powtcptest/internal/protocol/v1"
	"powtcptest/internal/tcp"
	"sync"
	"time"
)

var ErrFailedChallenge = errors.New("failed challenge")
var ErrInvalidChallenge = errors.New("invalid challenge")

type Server struct {
	addr                string
	readTimeout         time.Duration
	challengeDifficulty int

	listener          net.Listener
	collection        QuotesCollection
	challengeProvider PoWChallengeProvider

	connsMutex  sync.RWMutex
	activeConns map[string]*connInfo
	wg          sync.WaitGroup

	shutdownChan chan struct{}
}

func NewServer(
	listenAddr string,
	readTimeout time.Duration,
	challengeDifficulty int,
	collection QuotesCollection,
	challengeProvider PoWChallengeProvider,
) *Server {
	return &Server{
		addr:                listenAddr,
		readTimeout:         readTimeout,
		challengeDifficulty: challengeDifficulty,
		collection:          collection,
		challengeProvider:   challengeProvider,
		activeConns:         make(map[string]*connInfo),
		shutdownChan:        make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s: %w", s.addr, err)
	}
	s.listener = ln

	log.Printf("TCP server listening on port %s...", s.addr)

	go s.acceptConnections()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	close(s.shutdownChan)
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("error closing listener: %v", err)
	}

	// Wait for all active connections to finish or timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait() // Wait for all connections to complete
		close(done)
	}()

	select {
	case <-done:
		log.Println("All connections closed, server shutdown complete")
	case <-ctx.Done():
		log.Println("Shutdown timeout, forcefully closing remaining connections")
	}

	return nil
}

func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdownChan:
			log.Println("No longer accepting connections")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				var opErr *net.OpError
				if errors.As(err, &opErr) && opErr.Err.Error() == "use of closed network connection" {
					log.Println("Listener closed, stop accepting connections")
					return
				}
				log.Printf("Error accepting connection: %s", err.Error())
				continue
			}

			challenge := s.challengeProvider.GenerateChallenge()
			connInfo := newConnInfo(conn, challenge, s.challengeDifficulty)

			s.connsMutex.Lock()
			s.activeConns[conn.RemoteAddr().String()] = connInfo
			s.connsMutex.Unlock()

			log.Printf("Accepted new connection from %s", conn.RemoteAddr().String())
			s.wg.Add(1)
			go s.handleConnection(connInfo)
		}
	}
}

func (s *Server) closeConn(conn net.Conn) {
	s.connsMutex.Lock()
	defer s.connsMutex.Unlock()

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err.Error())
	}

	delete(s.activeConns, conn.RemoteAddr().String())

	log.Printf("Client %s disconnected", conn.RemoteAddr().String())
}

func (s *Server) handleConnection(connInfo *connInfo) {
	defer s.wg.Done()
	defer s.closeConn(connInfo.conn)

	newChallengeMsg := v1.NewNewChallengeMessage(connInfo.challenge, s.challengeDifficulty)
	if err := tcp.SendMessage(newChallengeMsg, connInfo.conn); err != nil {
		log.Printf("Error sending new challenge message: %s", err.Error())
		return
	}

	log.Printf("Successfully sent new challenge %s to %s", connInfo.challenge, connInfo.conn.RemoteAddr().String())

	for {
		select {
		case <-s.shutdownChan:
			log.Println("Server is shutting down. Closing connection...")
			return
		default:
			if err := connInfo.conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
				log.Printf("Error setting read deadline: %s", err.Error())
				return
			}

			// Read incoming message
			var msg v1.Message
			if err := json.NewDecoder(connInfo.conn).Decode(&msg); err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					continue
				}

				log.Printf("Error reading message: %v", err)
				return
			}

			err := s.handleClientMessage(msg, connInfo)
			switch {
			case errors.Is(err, ErrFailedChallenge):
				if _, err = connInfo.conn.Write([]byte("Challenge failed. Service denied.")); err != nil {
					log.Printf("Error writing failed challenge response: %v", err)
					return
				}
			case errors.Is(err, ErrInvalidChallenge):
				if _, err = connInfo.conn.Write([]byte("Challenge mismatch. Service denied.")); err != nil {
					log.Printf("Error writing failed challenge response: %v", err)
					return
				}
			case err != nil:
				log.Printf("Error handling client message: %v", err)
				return
			}

			log.Printf("Client %s served successfully", connInfo.conn.RemoteAddr().String())
			return
		}
	}
}

func (s *Server) handleClientMessage(msg v1.Message, connInfo *connInfo) error {
	switch msg.Type {
	case v1.MessageSolvedChallenge:
		solvedChallengeMsg, ok := msg.TypedMessage.(v1.SolvedChallengeMessage)
		if !ok {
			return fmt.Errorf("unable assert solved challenge message type")
		}
		//Check if it is the challenge we asked client to solve
		if solvedChallengeMsg.Challenge != connInfo.challenge {
			return ErrInvalidChallenge
		}

		if !s.challengeProvider.VerifyChallenge(solvedChallengeMsg.Challenge, solvedChallengeMsg.Solution, connInfo.challengeDifficulty) {
			return ErrFailedChallenge
		}

		quoteMsg := v1.NewWordOfWisdomMessage(s.collection.GetRandomQuote())
		if err := tcp.SendMessage(quoteMsg, connInfo.conn); err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("received message of unexpected type %d", msg.Type)
	}
}
