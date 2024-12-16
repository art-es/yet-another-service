//go:generate mockgen -source=retrier.go -destination=mock/retrier.go -package=mock
package retrier

import (
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/core/log"
)

type Retrier interface {
	Process(f func() error) error
}

type retrier struct {
	logger     log.Logger
	maxRetries int
	timeout    time.Duration
}

func New(logger log.Logger, maxRetries int, timeout time.Duration) Retrier {
	return &retrier{
		logger:     logger,
		maxRetries: maxRetries,
		timeout:    timeout,
	}
}

func (r *retrier) Process(f func() error) error {
	var retry int

	for {
		err := f()
		if err == nil {
			break
		}

		r.logger.Warn().Err(err).Msg("process failed in retrier")

		if retry == r.maxRetries {
			return fmt.Errorf("reached maximum number of retries: %w", err)
		}

		time.Sleep(r.timeout)

		retry++
	}

	return nil
}
