package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserStore struct {
	mock.Mock
}

func (mus *mockUserStore) FindUserByID(userID string) error {
	args := mus.Called(userID)
	return args.Error(0)
}

func TestAuthentication(t *testing.T) {
	mUserStore := new(mockUserStore)
	a := &Auth{mUserStore}
	secret := os.Getenv("JWT_SECRET_KEY")
	issuer := os.Getenv("JWT_NOTES_ISSUER")

	testCases := []struct {
		name        string
		userID      string
		bearerToken func() string
		setupMock   func()
		assertion   func(t *testing.T, err error)
	}{
		{
			name:        "Empty Bearer",
			bearerToken: func() string { return "" },
			setupMock:   func() {},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong format length",
			bearerToken: func() string { return "Bearer invalid length" },
			setupMock:   func() {},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong format Prefix",
			bearerToken: func() string { return "NoBearer asdf;lkj" },
			setupMock:   func() {},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: expected authorization header format: Bearer <token>")
			},
		},
		{
			name:        "Wrong method",
			bearerToken: func() string { return getBearerTokenEcdsa256(t) },
			setupMock:   func() {},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name:        "Invalid Token",
			bearerToken: func() string { return "Bearer invalidToken" },
			setupMock:   func() {},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "False secret",
			bearerToken: func() string {
				return setupJwtTokenString(t, jwt.MapClaims{}, "falseSecret")
			},
			setupMock: func() {},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "Expired Token",
			bearerToken: func() string {
				oneMinuteAgo := jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
				claims := setupClaims(oneMinuteAgo, "", "")
				return setupJwtTokenString(t, claims, secret)
			},
			setupMock: func() {},
			assertion: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "authenticate: error parsing tokenString")
			},
		},
		{
			name: "No expiration date set",
			bearerToken: func() string {
				return setupJwtTokenString(t, jwt.MapClaims{}, secret)
			},
			setupMock: func() {},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: no expiration date set")
			},
		},
		{
			name: "False issuer",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, "false issuer", "")
				return setupJwtTokenString(t, claims, secret)
			},
			setupMock: func() {},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: incorrect issuer")
			},
		},
		{
			name:   "User not found",
			userID: "000",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, issuer, "000")
				return setupJwtTokenString(t, claims, secret)
			},
			setupMock: func() {
				mUserStore.On("FindUserByID", "000").Return(errors.New("user not found"))
			},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: checkUserID: user not found")
				mUserStore.AssertCalled(t, "FindUserByID", "000")
			},
		},
		{
			name:   "userID (123) in endpoint unequal userID in jwt (456)",
			userID: "123",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, issuer, "456")
				return setupJwtTokenString(t, claims, secret)
			},
			setupMock: func() {
				mUserStore.On("FindUserByID", "123").Return(nil)
			},
			assertion: func(t *testing.T, err error) {
				assert.EqualError(t, err, "authenticate: user not enabled")
				mUserStore.AssertCalled(t, "FindUserByID", "123")
			},
		},
		{
			name:   "userID in endpoint unequal userID in jwt",
			userID: "123",
			bearerToken: func() string {
				inOneHour := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
				claims := setupClaims(inOneHour, issuer, "123")
				return setupJwtTokenString(t, claims, secret)
			},
			setupMock: func() {
				mUserStore.On("FindUserByID", "123").Return(nil)
			},
			assertion: func(t *testing.T, err error) {
				assert.NoError(t, err)
				mUserStore.AssertCalled(t, "FindUserByID", "123")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := a.Authenticate(tc.userID, tc.bearerToken())
			tc.assertion(t, err)
		})
	}
}

func ErrorContainss(t *testing.T, err error, containss ...string) {
	t.Helper()
	for _, contains := range containss {
		assert.ErrorContains(t, err, contains)
	}
}

func getBearerTokenEcdsa256(t *testing.T) (tokenString string) {
	t.Helper()
	var (
		key   *ecdsa.PrivateKey
		token *jwt.Token
	)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	token = jwt.New(jwt.SigningMethodES256)
	tokenString, err = token.SignedString(key)
	assert.NoError(t, err)
	return "Bearer " + tokenString
}

func newEmptyGetRequest(t *testing.T) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.NoError(t, err)
	return req
}

func addAuthorizationJWT(t *testing.T, tokenS string, req *http.Request) *http.Request {
	req.Header.Add("Authorization", "Bearer "+tokenS)
	return req
}

func setupJwtTokenString(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	bearerToken := "Bearer " + tokenS
	return bearerToken
}

func setupClaims(exp *jwt.NumericDate, iss, sub string) jwt.MapClaims {
	claims := jwt.MapClaims{
		"exp": exp,
		"iss": iss,
		"sub": sub,
	}
	return claims
}
