package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

func TestHashService(t *testing.T) {
	hashService := &HashService{}

	t.Run("ok", func(t *testing.T) {
		hashStr, err := hashService.Generate("foo")
		assert.NoError(t, err)
		assert.NotEmpty(t, hashStr)

		err = hashService.Check("foo", hashStr)
		assert.NoError(t, err)
	})

	t.Run("mismatched", func(t *testing.T) {
		hashStr, err := hashService.Generate("foo")
		assert.NoError(t, err)
		assert.NotEmpty(t, hashStr)

		err = hashService.Check("bar", hashStr)
		assert.ErrorIs(t, err, apperrors.ErrHashMismatched)
	})

	t.Run("wrong hash", func(t *testing.T) {
		err := hashService.Check("bar", "foo")
		assert.Error(t, err)
		assert.NotErrorIs(t, err, apperrors.ErrHashMismatched)
	})
}
