package client

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"powtcptest/internal/tcp/mock"
	"powtcptest/internal/tcp/server"
	"testing"
	"time"
)

func TestClientConnectionToServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	quotesCollection := mock.NewQuotesCollection(ctrl)
	challengeProvider := mock.NewPoWChallengeProvider(ctrl)
	srv := server.NewServer(":8001", 5*time.Second, 6, quotesCollection, challengeProvider)

	srv.Start()
	defer srv.Shutdown(context.Background())

	challengeSolver := mock.NewPoWChallengeSolver(ctrl)
	client := NewClient(":8001", 5*time.Second, challengeSolver)

	chg := "somechallenge"
	challengeProvider.EXPECT().GenerateChallenge().Return(chg)
	challengeProvider.EXPECT().VerifyChallenge(chg, gomock.Any(), gomock.Any()).Return(true)
	challengeSolver.EXPECT().SolveChallenge(gomock.Any(), chg, gomock.Any()).Return("solution", nil)
	quotesCollection.EXPECT().GetRandomQuote().Return("random_quote")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wisdom, err := client.RequestWordOfWisdom(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wisdom)
}

func TestClientConnectionTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewClient("127.0.0.1:9999", 1*time.Second, mock.NewPoWChallengeSolver(ctrl)) // Non-existent server

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wisdom, err := client.RequestWordOfWisdom(ctx)
	require.Error(t, err)
	require.Empty(t, wisdom)
}
