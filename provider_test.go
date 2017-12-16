package config

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentProvider(t *testing.T) {
	provider := NewEnvironmentProvider()
	assert.NotNil(t, provider)
	provider.lookup = func(_ string) (string, bool) {
		return "ok", true
	}

	err := provider.Configure(nil)
	assert.NoError(t, err)

	secret, ok := provider.GetSecret("hey")
	assert.Equal(t, "ok", secret)
	assert.True(t, ok)
}

func TestVaultProvider(t *testing.T) {
	s := httptest.NewServer(&testVaultHandler{})
	defer s.Close()

	provider := NewVaultProvider("user")
	assert.NotNil(t, provider)

	_, nok := provider.GetSecret("test")
	assert.False(t, nok)

	cfg := new(testConfig)
	err := provider.Configure(cfg)
	assert.Error(t, err)

	cfg.VaultCfg = &VaultConfig{
		Address:     s.URL,
		Application: "app",
	}

	err = provider.Configure(cfg)
	assert.NoError(t, err)

	_, ok := provider.GetSecret("test")
	assert.True(t, ok)

}

func TestVaultCache(t *testing.T) {
	s := httptest.NewServer(&testVaultHandler{})
	defer s.Close()

	cfg := new(testConfig)
	cfg.VaultCfg = &VaultConfig{
		Address:     s.URL,
		Application: "app",
	}
	cache, err := NewVaultCache("3497f684e809", cfg)
	assert.NotNil(t, cache)
	assert.NoError(t, err)
}
