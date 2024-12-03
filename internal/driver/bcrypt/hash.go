package bcrypt

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/art-es/yet-another-service/internal/domain/hash"
)

type HashService struct{}

func NewHashService() *HashService {
	return &HashService{}
}

func (s *HashService) Generate(str string) (string, error) {
	hashStr, err := bcrypt.GenerateFromPassword([]byte(str), 10)

	return string(hashStr), err
}

func (s *HashService) Check(str, hashStr string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(str))
	if err == nil {
		return nil
	}

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return hash.ErrMismatched
	}

	return err
}
