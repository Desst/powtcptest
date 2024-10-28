package server

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net"
	v1 "powtcptest/internal/protocol/v1"
	"powtcptest/internal/tcp"
	"powtcptest/internal/tcp/mock"
	"testing"
	"time"
)

func TestHandleMultipleConnections(t *testing.T) {
	ctrl := gomock.NewController(t)
	quotesCollection := mock.NewQuotesCollection(ctrl)
	challengeProvider := mock.NewPoWChallengeProvider(ctrl)
	srv := NewServer(":8001", 5*time.Second, 6, quotesCollection, challengeProvider)

	require.NoError(t, srv.Start()) // Start listening in the background

	challengeProvider.EXPECT().GenerateChallenge().Return("chg").Times(5)
	var conns []net.Conn
	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", srv.addr)
		require.NoError(t, err)
		conns = append(conns, conn)
	}

	time.Sleep(200 * time.Millisecond)
	err := srv.Shutdown(context.Background())
	require.NoError(t, err)

	require.Len(t, srv.activeConns, 0)
}

func TestInvalidChallenge(t *testing.T) {
	ctrl := gomock.NewController(t)
	quotesCollection := mock.NewQuotesCollection(ctrl)
	challengeProvider := mock.NewPoWChallengeProvider(ctrl)
	srv := NewServer(":8001", 5*time.Second, 6, quotesCollection, challengeProvider)

	require.NoError(t, srv.Start()) // Start listening in the background

	defer srv.Shutdown(context.Background())

	challengeProvider.EXPECT().GenerateChallenge().Return("some_challenge")

	conn, err := net.Dial("tcp", srv.addr)
	require.NoError(t, err)
	defer conn.Close()

	// Read incoming challenge message
	var msg v1.Message
	var newChallengeMsg v1.NewChallengeMessage
	require.NoError(t, json.NewDecoder(conn).Decode(&msg))
	newChallengeMsg, ok := msg.TypedMessage.(v1.NewChallengeMessage)
	require.True(t, ok)
	require.NotEmpty(t, newChallengeMsg.Challenge)

	// Send an invalid challenge
	err = tcp.SendMessage(v1.NewSolvedChallengeMessage("some_other_challenge", "some_solution"), conn)
	require.NoError(t, err)

	// Expect the server to send an error msg
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	require.NoError(t, err)
	require.Equal(t, string(buf[:n]), "Challenge mismatch. Service denied.")
}

func TestShutdownWhileConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	quotesCollection := mock.NewQuotesCollection(ctrl)
	challengeProvider := mock.NewPoWChallengeProvider(ctrl)
	srv := NewServer(":8001", 5*time.Second, 6, quotesCollection, challengeProvider)
	require.NoError(t, srv.Start()) // Start listening in the background

	challengeProvider.EXPECT().GenerateChallenge().Return("some_challenge")
	conn, err := net.Dial("tcp", srv.addr)
	require.NoError(t, err)

	//expect newchallenge msg
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	require.NoError(t, err) // Connection should be closed during shutdown

	// Trigger a server shutdown while the connection is open
	require.NoError(t, srv.Shutdown(context.Background()))

	n, err := conn.Read(buf)
	require.Error(t, err) // Connection should be closed during shutdown
	require.Equal(t, 0, n)

	conn.Close()
}
