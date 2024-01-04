package fuego

import (
	"net/http"
	"net/http/httptest"
	"regexp"
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

func TestCheckRolesOr(t *testing.T) {
	check := checkRolesOr("a", "b")

	t.Run("empty", func(t *testing.T) {
		require.False(t, check())
	})

	t.Run("one", func(t *testing.T) {
		require.True(t, check("a"))
		require.True(t, check("b"))
		require.False(t, check("c"))
	})

	t.Run("multiple", func(t *testing.T) {
		require.True(t, check("a", "b"))
		require.True(t, check("a", "c"))
		require.True(t, check("b", "c"))
		require.False(t, check("c", "d"))
	})

	t.Run("nobody accepted", func(t *testing.T) {
		check := checkRolesOr()
		require.False(t, check())
		require.False(t, check("a", "c"))
		require.False(t, check("b", "c"))
	})
}

func TestCheckRolesRegex(t *testing.T) {
	re := regexp.MustCompile(`^a.*`)
	check := checkRolesRegex(re)

	t.Run("empty", func(t *testing.T) {
		require.False(t, check())
	})

	t.Run("one", func(t *testing.T) {
		require.True(t, check("a"))
		require.True(t, check("ab"))
		require.False(t, check("b"))
	})

	t.Run("multiple", func(t *testing.T) {
		require.True(t, check("a", "b"))
		require.True(t, check("a", "c"))
		require.False(t, check("b", "c"))
		require.False(t, check("ca", "d"))
	})
}

func TestTokenFromHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer 123")

	token := TokenFromHeader(r)
	require.Equal(t, "123", token)
}

func TestTokenFromCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: JWTCookieName, Value: "456"})

	token := TokenFromCookie(r)
	require.Equal(t, "456", token)
}
