package layers

import (
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/jitter"
	"github.com/Rican7/retry/strategy"
	"math/rand"
	"time"
)

func ExampleNewRetryLayer() {
	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	NewRetryLayer(
		SetStrategy(
			strategy.Limit(5),
			strategy.BackoffWithJitter(
				backoff.BinaryExponential(10*time.Millisecond),
				jitter.Deviation(random, 0.5),
			),
		),
	)
}
