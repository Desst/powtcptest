//go:generate go run github.com/golang/mock/mockgen -destination=../mock/pow_challenge_solver.go -package=mock -mock_names PoWChallengeSolver=PoWChallengeSolver powtcptest/internal/tcp/client PoWChallengeSolver

package client

import "context"

type PoWChallengeSolver interface {
	SolveChallenge(ctx context.Context, challenge string, difficulty int) (string, error)
}
