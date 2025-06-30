package fuego

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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

	t.Run("can initialize with custom key", func(t *testing.T) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)
		require.NotEmpty(t, key)

		security := NewSecurityWithKey(key)
		require.NotEmpty(t, security)

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
		fmt.Printf("error: %v\n", err)
		require.ErrorAs(t, err, &UnauthorizedError{})
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

func TestTokenFromQueryParam(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)

		token := TokenFromQueryParam(r)
		require.Empty(t, token)
	})

	t.Run("with token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/?jwt=123", nil)

		token := TokenFromQueryParam(r)
		require.Equal(t, "123", token)
	})
}

func TestTokenFromHeader(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)

		token := TokenFromHeader(r)
		require.Empty(t, token)
	})

	t.Run("with invalid token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", "Bla")

		token := TokenFromHeader(r)
		require.Empty(t, token)
	})

	t.Run("with invalid token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", "Blabla 123")

		token := TokenFromHeader(r)
		require.Empty(t, token)
	})

	t.Run("with valid token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", "Bearer 123")

		token := TokenFromHeader(r)
		require.Equal(t, "123", token)
	})
}

func TestTokenFromCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: JWTCookieName, Value: "456"})

	token := TokenFromCookie(r)
	require.Equal(t, "456", token)
}

func TestAuthWall(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("list", func(t *testing.T) {
		authWall := AuthWall("a", "b")
		require.NotNil(t, authWall)
		t.Run("no token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{
				"sub":   "123",
				"roles": []string{"c", "d"},
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			authWall(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusForbidden, w.Code)
		})

		t.Run("with token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{
				"sub":   "123",
				"roles": []string{"a", "d"},
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			authWall(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("regex", func(t *testing.T) {
		authWall := AuthWallRegex(`^a.*`)
		require.NotNil(t, authWall)
		t.Run("no token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{
				"sub":   "123",
				"roles": []string{"c", "d"},
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			authWall(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusForbidden, w.Code)
		})

		t.Run("with token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{
				"sub":   "123",
				"roles": []string{"a", "d"},
			})
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()

			authWall(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusOK, w.Code)
		})
	})
}

func TestTokenFromContext(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		token, err := TokenFromContext(r.Context())
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("with invalid token type", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(r.Context(), contextKeyJWT, "123")
		token, err := TokenFromContext(ctx)
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("with token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{"sub": "123"})
		token, err := TokenFromContext(ctx)
		require.NoError(t, err)
		sub, err := token.GetSubject()
		require.NoError(t, err)
		require.Equal(t, "123", sub)
	})
}

func TestGenerateTokenToCookies(t *testing.T) {
	security := NewSecurity()
	claims := jwt.MapClaims{
		"aud": "test",
		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
		"iss": "test",
		"nbf": jwt.NewNumericDate(time.Now()),
		"sub": "123",
	}

	w := httptest.NewRecorder()
	security.GenerateTokenToCookies(claims, w)

	authCookie := w.Result().Cookies()[0]
	require.NotEmpty(t, authCookie)
	require.Equal(t, JWTCookieName, authCookie.Name)
}

func TestTokenToContext(t *testing.T) {
	security := NewSecurity()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("from header", func(t *testing.T) {
		tokenToContext := security.TokenToContext(
			TokenFromHeader,
		)

		t.Run("no token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			tokenToContext(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("wrong token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Authorization", "Bearer 123")
			w := httptest.NewRecorder()
			tokenToContext(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("correct token", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			token, err := security.GenerateToken(jwt.MapClaims{"sub": "123"})
			require.NoError(t, err)

			r.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			tokenToContext(h).ServeHTTP(w, r)
			require.Equal(t, http.StatusOK, w.Code)
		})
	})
}

func TestSecurity_CookieLogoutHandler(t *testing.T) {
	security := NewSecurity()

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	security.CookieLogoutHandler(w, r)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	authCookie := cookies[0]
	require.NotEmpty(t, authCookie)
	require.Equal(t, JWTCookieName, authCookie.Name)
}

func TestSecurity_RefreshHandler(t *testing.T) {
	security := NewSecurity()

	t.Run("no token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		security.RefreshHandler(w, r)

		cookies := w.Result().Cookies()
		require.Empty(t, cookies)
	})

	t.Run("with token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		ctx := WithValue(r.Context(), jwt.MapClaims{
			"aud": "test",
			"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			"iat": jwt.NewNumericDate(time.Now()),
			"iss": "test",
			"nbf": jwt.NewNumericDate(time.Now()),
			"sub": "123",
		})
		r = r.WithContext(ctx)

		security.RefreshHandler(w, r)

		body := w.Body.String()
		t.Log(body)
		require.Equal(t, http.StatusOK, w.Code)
		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		authCookie := cookies[0]
		require.Equal(t, JWTCookieName, authCookie.Name)
	})
}

func TestSecurity_StdLoginHandler(t *testing.T) {
	security := NewSecurity()
	v := func(r *http.Request) (jwt.Claims, error) {
		if r.FormValue("user") != "test" || r.FormValue("password") != "test" {
			return nil, UnauthorizedError{}
		}
		return jwt.MapClaims{"sub": "123"}, nil
	}
	loginHandler := security.StdLoginHandler(v)

	t.Run("with incorrect ids", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		loginHandler(w, r)

		cookies := w.Result().Cookies()
		require.Empty(t, cookies)
	})

	t.Run("with correct ids", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/?user=test&password=test", nil)
		w := httptest.NewRecorder()
		loginHandler(w, r)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		authCookie := cookies[0]
		require.NotEmpty(t, authCookie)
		require.Equal(t, JWTCookieName, authCookie.Name)
	})
}

func TestSecurity_LoginHandler(t *testing.T) {
	security := NewSecurity()
	v := func(user, password string) (jwt.Claims, error) {
		if user != "test" || password != "test" {
			return nil, UnauthorizedError{}
		}
		return jwt.MapClaims{"sub": "123"}, nil
	}
	loginHandler := security.LoginHandler(v)

	t.Run("without ids", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		s := NewServer()
		truc := HTTPHandler(s, loginHandler, BaseRoute{})
		truc.ServeHTTP(w, r)

		cookies := w.Result().Cookies()
		require.Empty(t, cookies)
	})

	t.Run("with incorrect ids", func(t *testing.T) {
		loginBody := `{"user": "hacker", "password": "hacker"}`
		r := httptest.NewRequest("GET", "/", strings.NewReader(loginBody))
		w := httptest.NewRecorder()

		s := NewServer()
		truc := HTTPHandler(s, loginHandler, BaseRoute{})
		truc.ServeHTTP(w, r)

		cookies := w.Result().Cookies()
		require.Empty(t, cookies)
	})

	t.Run("with correct ids", func(t *testing.T) {
		loginBody := `{"user": "test", "password": "test"}`
		r := httptest.NewRequest("GET", "/", strings.NewReader(loginBody))
		w := httptest.NewRecorder()

		s := NewServer()
		truc := HTTPHandler(s, loginHandler, BaseRoute{})
		truc.ServeHTTP(w, r)

		cookies := w.Result().Cookies()
		require.Len(t, cookies, 1)
		authCookie := cookies[0]
		require.NotEmpty(t, authCookie)
		require.Equal(t, JWTCookieName, authCookie.Name)
	})
}

func TestGetToken(t *testing.T) {
	t.Run("no token", func(t *testing.T) {
		ctx := context.Background()
		token, err := GetToken[any](ctx)
		require.Error(t, err)
		require.Empty(t, token)
	})

	t.Run("with valid token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(r.Context(), contextKeyJWT, jwt.MapClaims{"sub": "123"})

		token, err := GetToken[jwt.MapClaims](ctx)
		require.NoError(t, err)
		sub, err := token.GetSubject()
		require.NoError(t, err)
		require.Equal(t, "123", sub)
	})

	t.Run("with token of custom type", func(t *testing.T) {
		type MyToken struct {
			jwt.MapClaims
			Username string
			UserID   string
		}
		r := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(r.Context(), contextKeyJWT, MyToken{MapClaims: jwt.MapClaims{"sub": "123"}})

		_, err := GetToken[MyToken](ctx)
		require.Error(t, err)
	})
}
