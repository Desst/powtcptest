package pow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateChallenge() string {
	return fmt.Sprintf("%x", rand.Int63())
}

// VerifyChallenge - checks if the provided solution has <difficulty> number of leading zeros
func (s *Service) VerifyChallenge(challenge string, solution string, difficulty int) bool {
	hash := sha256.New()
	hash.Write([]byte(challenge + solution))
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	return strings.HasPrefix(hashString, strings.Repeat("0", difficulty))
}

func (s *Service) SolveChallenge(ctx context.Context, challenge string, difficulty int) (string, error) {
	var solution int
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			// Try different solutions
			hash := sha256.New()
			hash.Write([]byte(challenge + fmt.Sprintf("%d", solution)))
			hashSum := hash.Sum(nil)
			hashString := hex.EncodeToString(hashSum)
			if strings.HasPrefix(hashString, strings.Repeat("0", difficulty)) {
				return fmt.Sprintf("%d", solution), nil
			}
			solution++
		}
	}
}
