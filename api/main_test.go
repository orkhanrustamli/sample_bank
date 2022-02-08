package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mockdb "github.com/orkhanrustamli/simplebank/db/mock"
	"github.com/orkhanrustamli/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createtestServer(t *testing.T, store *mockdb.MockStore) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration:     time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
