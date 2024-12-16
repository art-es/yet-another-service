package retrier

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/art-es/yet-another-service/internal/driver/zerolog"
)

func TestRetrier(t *testing.T) {
	for _, tt := range []struct {
		name       string
		newProcess func(t *testing.T) func() error
		assert     func(t *testing.T, err error, logs []string)
	}{
		{
			name: "no retries",
			newProcess: func(t *testing.T) func() error {
				var count int
				return func() error {
					count++
					assert.Equal(t, 1, count)
					return nil
				}
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Len(t, logs, 0)
			},
		},
		{
			name: "retried",
			newProcess: func(t *testing.T) func() error {
				var count int
				return func() error {
					count++
					assert.LessOrEqual(t, count, 3)
					if count < 3 {
						return fmt.Errorf("error#%d", count)
					}
					return nil
				}
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.NoError(t, err)
				assert.Len(t, logs, 2)
				assert.Equal(t, `{"level":"warn","error":"error#1","message":"process failed in retrier"}`, logs[0])
				assert.Equal(t, `{"level":"warn","error":"error#2","message":"process failed in retrier"}`, logs[1])
			},
		},
		{
			name: "reached max retries",
			newProcess: func(t *testing.T) func() error {
				var count int
				return func() error {
					count++
					assert.LessOrEqual(t, count, 5)
					return fmt.Errorf("error#%d", count)
				}
			},
			assert: func(t *testing.T, err error, logs []string) {
				assert.EqualError(t, err, "reached maximum number of retries: error#5")
				assert.Len(t, logs, 5)
				assert.Equal(t, `{"level":"warn","error":"error#1","message":"process failed in retrier"}`, logs[0])
				assert.Equal(t, `{"level":"warn","error":"error#2","message":"process failed in retrier"}`, logs[1])
				assert.Equal(t, `{"level":"warn","error":"error#3","message":"process failed in retrier"}`, logs[2])
				assert.Equal(t, `{"level":"warn","error":"error#4","message":"process failed in retrier"}`, logs[3])
				assert.Equal(t, `{"level":"warn","error":"error#5","message":"process failed in retrier"}`, logs[4])
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			logbuf := &bytes.Buffer{}
			logger := zerolog.NewLoggerWithWriter(logbuf)
			retrier := New(logger, 4, 0)
			process := tt.newProcess(t)

			err := retrier.Process(process)
			logs := logsFromBuffer(logbuf)

			tt.assert(t, err, logs)
		})
	}
}

func logsFromBuffer(buf *bytes.Buffer) []string {
	var out []string
	for s := bufio.NewScanner(buf); s.Scan(); {
		out = append(out, s.Text())
	}
	return out
}
