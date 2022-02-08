package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orkhanrustamli/simplebank/token"
	"github.com/stretchr/testify/require"
)

const (
	username = "someUsername"
	duration = time.Minute
)

func setAuthHeader(
	request *http.Request,
	authHeader string,
	authHeaderType string,
	token string,
) {
	h := fmt.Sprintf("%s %s", authHeaderType, token)

	request.Header.Set(authHeader, h)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenManager token.Manager)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.Manager) {
				token, err := tokenManager.CreateToken(username, duration)
				require.NoError(t, err)

				setAuthHeader(request, authorizationHeader, authorizationHeaderType, token)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "HeaderNotProvided",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.Manager) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "WrongType",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.Manager) {
				token, err := tokenManager.CreateToken(username, duration)
				require.NoError(t, err)

				setAuthHeader(request, authorizationHeader, "jwt", token)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.Manager) {
				token, err := tokenManager.CreateToken(username, -duration)
				require.NoError(t, err)

				setAuthHeader(request, authorizationHeader, authorizationHeaderType, token)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := createtestServer(t, nil)

			path := "/auth"
			server.router.GET(path, authMiddleware(server.tokenManager), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			request, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
