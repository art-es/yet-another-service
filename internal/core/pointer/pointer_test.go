package pointer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTo(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		res := To((*int)(nil))
		assert.Equal(t, "**int", fmt.Sprintf("%T", res))
		assert.NotNil(t, res)
		assert.Nil(t, *res)
	})

	t.Run("value", func(t *testing.T) {
		res := To(1)
		assert.NotNil(t, res)
		assert.Equal(t, 1, *res)
	})
}

func TestFrom(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		ptr := (*int)(nil)
		res := From(ptr)
		assert.Equal(t, 0, res)
	})

	t.Run("not nil pointer", func(t *testing.T) {
		val := 1
		res := From(&val)
		assert.Equal(t, val, res)
	})
}
