package pow

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChallenges(t *testing.T) {
	t.Parallel()

	powService := NewService()

	for i := 0; i < 8; i++ {
		difficulty := i
		t.Run(fmt.Sprintf("Difficulty %d", i), func(t *testing.T) {
			challenge := powService.GenerateChallenge()
			start := time.Now()
			solution, err := powService.SolveChallenge(context.Background(), challenge, difficulty)
			elapsed := time.Since(start)
			require.NoError(t, err)
			log.Printf("Challenge with difficulty %d solved in %.2f seconds", difficulty, elapsed.Seconds())
			require.True(t, powService.VerifyChallenge(challenge, solution, difficulty))
		})
	}
}
