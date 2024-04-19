package common

import "crypto/rand"

func MustGenerateRandomKey(keyLength int) []byte {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}

// key, err := generateRandomKey(32)
// assert.NoError(t, err)
// jwtSvc := jwt.NewJWT(key)

// testCases := []struct {
// 	name       string
// 	tokenS     func() string
// 	assertions func(t *testing.T, err error)
// }{
// 	{
// 		name:   "Invalid Token",
// 		tokenS: func() string { return "invalidToken" },
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: error parsing tokenString")
// 		},
// 	},
// 	{
// 		name: "False secret",
// 		tokenS: func() string {
// 			token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
// 			tokenS, err := token.SignedString([]byte("false secret"))
// 			assert.NoError(t, err)
// 			return tokenS
// 		},
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: error parsing tokenString")
// 		},
// 	},
// 	{
// 		name: "Expired Token",
// 		tokenS: func() string {
// 			oneMinuteAgo := jwtLib.NewNumericDate(time.Now().Add(-1 * time.Minute))
// 			claims := setupClaims(oneMinuteAgo, "", "")
// 			return setupJwtTokenString(t, claims, j.key)
// 		},
// 		assertions: func(t *testing.T, err error) {
// 			assert.ErrorContains(t, err, "verify: ")
// 		},
// 	},
// 	// {
// 	// 	name: "No expiration date set",
// 	// 	tokenS: func() string {

// 	// 		token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, jwtLib.MapClaims{})
// 	// 		tokenS, err := token.SignedString([]byte(j.key))
// 	// 		assert.NoError(t, err)
// 	// 		return tokenS
// 	// 	},
// 	// 	assertions: func(t *testing.T, err error) {
// 	// 		assert.EqualError(t, err, "authenticate: no expiration date set")
// 	// 	},
// 	// },
// 	// {
// 	// 	name: "False issuer",
// 	// 	bearerToken: func() string {
// 	// 		inOneHour := jwtLib.NewNumericDate(time.Now().Add(1 * time.Hour))
// 	// 		claims := setupClaims(inOneHour, "false issuer", "")
// 	// 		return setupJwtTokenString(t, claims, secret)
// 	// 	},
// 	// 	assertion: func(t *testing.T, err error) {
// 	// 		assert.EqualError(t, err, "authenticate: incorrect issuer")
// 	// 	},
// 	// },
// 	//
// 	// {
// 	// 	name:        "Wrong method",
// 	// 	bearerToken: func() string { return getBearerTokenEcdsa256(t) },
// 	// 	assertion: func(t *testing.T, err error) {
// 	// 		assert.ErrorContains(t, err, "authenticate: verify: ")
// 	// 	},
// 	// },
// }

// for _, tc := range testCases {
// 	_, err := j.Verify(tc.tokenS())
// 	tc.assertions(t, err)
// }
//
// func setupJwtTokenString(t *testing.T, claims jwtLib.MapClaims, secret string) string {
// 	t.Helper()
// 	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
// 	tokenS, err := token.SignedString([]byte(secret))
// 	assert.NoError(t, err)
// 	bearerToken := "Bearer " + tokenS
// 	return bearerToken
// }
