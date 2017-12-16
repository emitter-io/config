package config

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testVaultHandler struct{}

func (h *testVaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := ""
	switch r.URL.Path {
	case "/v1/auth/app-id/login":
		response = `{ "auth": {"client_token": "123" } }`
	case "/v1/secret/test":
		response = `{ "data": {"value": "hi" } }`
	case "/v1/secret/err1":
		response = `err`
	case "/v1/secret/err2":
		response = `{ }`
	case "/v1/aws/sts/storage-emitter":
		response = `{ "lease_duration": 100, "data": { "access_key": "access", "secret_key": "secret", "security_token": "token" } }`
	case "/v1/aws/sts/err1":
		response = `err`
	case "/v1/aws/sts/err2":
		response = `{ }`
	}

	w.Write([]byte(response))
}

func Test_newVaultClient(t *testing.T) {
	tests := []struct {
		addr string
		err  bool
	}{
		{addr: "127.0.0.1"},
		{addr: "127.0.0.1:8080"},
		{addr: "http://127.0.0.1:8080"},
	}

	for _, tc := range tests {
		cli := NewVaultClient(tc.addr)
		assert.NotNil(t, cli)
	}
}

func TestVaultClient(t *testing.T) {
	s := httptest.NewServer(&testVaultHandler{})
	defer s.Close()

	// Test authentication first
	cli := NewVaultClient(s.URL)
	err := cli.Authenticate("xxxxxx", "yyyyyy")
	assert.NoError(t, err)
	assert.Equal(t, "123", cli.token)

	// Test few different secret endpoints
	secretTests := []struct {
		key string
		err bool
	}{
		{key: "test"},
		{key: "err1", err: true},
		{key: "err2", err: true},
	}

	for _, tc := range secretTests {
		v, err := cli.ReadSecret(tc.key)
		assert.Equal(t, tc.err, err != nil)
		if !tc.err {
			assert.Equal(t, "hi", v)
		}
	}

}
