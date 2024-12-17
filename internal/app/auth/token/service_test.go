package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/auth/token/mock"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

func TestGenerate(t *testing.T) {
	type mocks struct {
		jwtService *mock.MockjwtService
	}

	now, _ := time.Parse(time.DateTime, "2000-01-01 10:00:00")
	nextHour, _ := time.Parse(time.DateTime, "2000-01-01 11:00:00")
	nextWeek, _ := time.Parse(time.DateTime, "2000-01-08 10:00:00")

	getCurrentTime = func() time.Time {
		return now
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, res *auth.TokenPair, err error)
	}{
		{
			name: "generate access token error",
			setup: func(t *testing.T, m mocks) {
				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.jwtService.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.TokenPair, err error) {
				assert.EqualError(t, err, "generate access token: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "generate refresh token error",
			setup: func(t *testing.T, m mocks) {
				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.jwtService.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("dummy access token", nil)

				expRefreshTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextWeek,
					UserID:    "dummy user id",
				}

				m.jwtService.EXPECT().
					Generate(gomock.Eq(expRefreshTokenClaims)).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *auth.TokenPair, err error) {
				assert.EqualError(t, err, "generate refresh token: dummy error")
				assert.Nil(t, res)
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				expAccessTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextHour,
					UserID:    "dummy user id",
				}

				m.jwtService.EXPECT().
					Generate(gomock.Eq(expAccessTokenClaims)).
					Return("dummy access token", nil)

				expRefreshTokenClaims := &auth.TokenClaims{
					IssuedAt:  now,
					ExpiresAt: nextWeek,
					UserID:    "dummy user id",
				}

				m.jwtService.EXPECT().
					Generate(gomock.Eq(expRefreshTokenClaims)).
					Return("dummy refresh token", nil)
			},
			assert: func(t *testing.T, res *auth.TokenPair, err error) {
				assert.NoError(t, err)
				expResult := &auth.TokenPair{
					AccessToken:  "dummy access token",
					RefreshToken: "dummy refresh token",
				}
				assert.Equal(t, expResult, res)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				jwtService: mock.NewMockjwtService(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.jwtService, nil)
			res, err := service.Generate("dummy user id")

			if tt.assert != nil {
				tt.assert(t, res, err)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	type mocks struct {
		jwtService *mock.MockjwtService
		blackList  *mock.MockblackList
	}

	now, _ := time.Parse(time.DateTime, "2000-01-01 10:00:00")
	nextHour, _ := time.Parse(time.DateTime, "2000-01-01 11:00:00")

	getCurrentTime = func() time.Time {
		return now
	}

	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, m mocks)
		assert func(t *testing.T, accessToken string, err error)
	}{
		{
			name: "parse refresh token error",
			setup: func(t *testing.T, m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.EqualError(t, err, "parse refresh token: dummy error")
			},
		},
		{
			name: "check refresh token in black list error",
			setup: func(t *testing.T, m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy refresh token")).
					Return(false, errors.New("dummy error"))
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.EqualError(t, err, "check refresh token in black list: dummy error")
			},
		},
		{
			name: "refresh token in black list",
			setup: func(t *testing.T, m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy refresh token")).
					Return(true, nil)
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.ErrorIs(t, err, apperrors.ErrInvalidAuthToken)
			},
		},
		{
			name: "generate access token error",
			setup: func(t *testing.T, m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy refresh token")).
					Return(false, nil)

				m.jwtService.EXPECT().
					Generate(gomock.Eq(&auth.TokenClaims{
						IssuedAt:  now,
						ExpiresAt: nextHour,
						UserID:    "dummy user id",
					})).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.EqualError(t, err, "generate access token: dummy error")
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy refresh token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy refresh token")).
					Return(false, nil)

				m.jwtService.EXPECT().
					Generate(gomock.Eq(&auth.TokenClaims{
						IssuedAt:  now,
						ExpiresAt: nextHour,
						UserID:    "dummy user id",
					})).
					Return("dummy access token", nil)
			},
			assert: func(t *testing.T, accessToken string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "dummy access token", accessToken)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				jwtService: mock.NewMockjwtService(ctrl),
				blackList:  mock.NewMockblackList(ctrl),
			}

			if tt.setup != nil {
				tt.setup(t, m)
			}

			service := NewService(m.jwtService, m.blackList)
			accessToken, err := service.Refresh(context.Background(), "dummy refresh token")

			if tt.assert != nil {
				tt.assert(t, accessToken, err)
			}
		})
	}
}

func TestAuthorize(t *testing.T) {
	type mocks struct {
		jwtService *mock.MockjwtService
		blackList  *mock.MockblackList
	}

	for _, tt := range []struct {
		name   string
		setup  func(m mocks)
		assert func(t *testing.T, userID string, err error)
	}{
		{
			name: "parse access token error",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse("dummy access token").
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, userID string, err error) {
				assert.EqualError(t, err, "parse access token: dummy error")
				assert.Empty(t, userID)
			},
		},
		{
			name: "check access token in black list error",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy access token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy access token")).
					Return(false, errors.New("dummy error"))
			},
			assert: func(t *testing.T, userID string, err error) {
				assert.EqualError(t, err, "check access token in black list: dummy error")
				assert.Empty(t, userID)
			},
		},
		{
			name: "access token in black list",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy access token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy access token")).
					Return(true, nil)
			},
			assert: func(t *testing.T, userID string, err error) {
				assert.ErrorIs(t, err, apperrors.ErrInvalidAuthToken)
				assert.Empty(t, userID)
			},
		},
		{
			name: "ok",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy access token")).
					Return(&auth.TokenClaims{UserID: "dummy user id"}, nil)

				m.blackList.EXPECT().
					Has(gomock.Any(), gomock.Eq("dummy access token")).
					Return(false, nil)
			},
			assert: func(t *testing.T, userID string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "dummy user id", userID)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				jwtService: mock.NewMockjwtService(ctrl),
				blackList:  mock.NewMockblackList(ctrl),
			}

			tt.setup(m)

			service := NewService(m.jwtService, m.blackList)
			userID, err := service.Authorize(context.Background(), "dummy access token")

			tt.assert(t, userID, err)
		})
	}
}

func TestInvalidate(t *testing.T) {
	type mocks struct {
		jwtService *mock.MockjwtService
		blackList  *mock.MockblackList
	}

	now, _ := time.Parse(time.DateTime, "2000-01-01 11:00:00")
	prevHour, _ := time.Parse(time.DateTime, "2000-01-01 10:00:00")
	nextHour, _ := time.Parse(time.DateTime, "2000-01-01 12:00:00")

	getCurrentTime = func() time.Time {
		return now
	}

	for _, tt := range []struct {
		name   string
		setup  func(m mocks)
		assert func(t *testing.T, err error)
	}{
		{
			name: "parse token error",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy token")).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "parse token: dummy error")
			},
		},
		{
			name: "black list TTL is negative",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy token")).
					Return(&auth.TokenClaims{ExpiresAt: nextHour}, nil)
			},
			assert: func(t *testing.T, err error) {
				assert.EqualError(t, err, "black list TTL is negative")
			},
		},
		{
			name: "ok",
			setup: func(m mocks) {
				m.jwtService.EXPECT().
					Parse(gomock.Eq("dummy token")).
					Return(&auth.TokenClaims{ExpiresAt: prevHour}, nil)

				m.blackList.EXPECT().
					Add(gomock.Any(), gomock.Eq("dummy token"), gomock.Eq(time.Hour)).
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
				jwtService: mock.NewMockjwtService(ctrl),
				blackList:  mock.NewMockblackList(ctrl),
			}

			tt.setup(m)

			service := NewService(m.jwtService, m.blackList)
			err := service.Invalidate(context.Background(), "dummy token")

			tt.assert(t, err)
		})
	}
}
