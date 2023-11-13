package op

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type SecurityInfo struct {
	jwt.RegisteredClaims
	Username string
	UserID   string
}

func TestSecurity(t *testing.T) {
	now := time.Now()

	me := SecurityInfo{}
	me.IssuedAt = jwt.NewNumericDate(now)

	t.Run("working JWT encoding/decoding", func(t *testing.T) {
		security := NewSecurity()

		security.Now = func() time.Time { return now }
		security.ExpiresInterval = 10 * time.Minute

		s, err := security.GenerateToken(me)
		require.NoError(t, err)
		require.NotEmpty(t, s)

		security.Now = func() time.Time { return now.Add(5 * time.Minute) }
		decoded, err := security.ValidateToken(s)
		require.NoError(t, err)
		require.NotEmpty(t, decoded)
	})

	t.Run("expired", func(t *testing.T) {
		security := NewSecurity()

		security.Now = func() time.Time { return now }
		security.ExpiresInterval = 10 * time.Minute

		s, err := security.GenerateToken(me)
		require.NoError(t, err)
		require.NotEmpty(t, s)

		security.Now = func() time.Time { return now.Add(15 * time.Minute) }
		decoded, err := security.ValidateToken(s)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrExpired)
		require.Empty(t, decoded)
	})
}
