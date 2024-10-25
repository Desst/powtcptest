//go:generate go run github.com/golang/mock/mockgen -destination=../mock/quotes_collection.go -package=mock -mock_names QuotesCollection=QuotesCollection powtcptest/internal/tcp/server QuotesCollection
//go:generate go run github.com/golang/mock/mockgen -destination=../mock/pow_challenge_provider.go -package=mock -mock_names PoWChallengeProvider=PoWChallengeProvider powtcptest/internal/tcp/server PoWChallengeProvider

package server

type QuotesCollection interface {
	GetRandomQuote() string
}

type PoWChallengeProvider interface {
	GenerateChallenge() string
	VerifyChallenge(challenge string, solution string, difficulty int) bool
}
