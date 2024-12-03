package activate

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/domain/auth/activate/mock"
)

func TestService(t *testing.T) {
	type mocks struct {
		activationRepository *mock.MockactivationRepository
		userRepository       *mock.MockuserRepository
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, err error)
	}{
		{
			name: "find activation by token in repository error",
			setup: func(t *testing.T, m mocks) {
				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "find activation by token in repository: dummy error")
			},
		},
		{
			name: "activation not found",
			setup: func(t *testing.T, m mocks) {
				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(nil, nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, auth.ErrActivationNotFound)
			},
		},
		{
			name: "activate user in repository error",
			setup: func(t *testing.T, m mocks) {
				activation := &auth.Activation{
					Token:  "dummy token",
					UserID: "dummy user id",
				}

				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(activation, nil)

				m.userRepository.EXPECT().
					Activate(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "activate user in repository: dummy error")
			},
		},
		{
			name: "delete activation by token in repository error",
			setup: func(t *testing.T, m mocks) {
				activation := &auth.Activation{
					Token:  "dummy token",
					UserID: "dummy user id",
				}

				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(activation, nil)

				m.userRepository.EXPECT().
					Activate(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(nil)

				m.activationRepository.EXPECT().
					Delete(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy token")).
					Return(errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "delete activation by token in repository: dummy error")
			},
		},
		{
			name: "commit transaction error",
			setup: func(t *testing.T, m mocks) {
				activation := &auth.Activation{
					Token:  "dummy token",
					UserID: "dummy user id",
				}

				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(activation, nil)

				m.userRepository.EXPECT().
					Activate(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(nil)

				m.activationRepository.EXPECT().
					Delete(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy token")).
					Do(func(_ context.Context, tx transaction.Transaction, _ string) {
						tx.AddCommit(func() error {
							return errors.New("dummy error")
						})
					}).
					Return(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "commit transaction: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				activation := &auth.Activation{
					Token:  "dummy token",
					UserID: "dummy user id",
				}

				m.activationRepository.EXPECT().
					FindByToken(gomock.Any(), gomock.Eq("dummy token")).
					Return(activation, nil)

				m.userRepository.EXPECT().
					Activate(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy user id")).
					Return(nil)

				m.activationRepository.EXPECT().
					Delete(gomock.Any(), gomock.Not(nil), gomock.Eq("dummy token")).
					Return(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				activationRepository: mock.NewMockactivationRepository(ctrl),
				userRepository:       mock.NewMockuserRepository(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.activationRepository, m.userRepository)
			err := service.Activate(context.Background(), "dummy token")

			if tt.assert != nil {
				tt.assert(t, err)
			}
		})
	}
}
