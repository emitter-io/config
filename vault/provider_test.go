// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package vault

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaultProvider(t *testing.T) {
	s := httptest.NewServer(&testVaultHandler{})
	defer s.Close()

	provider := NewProvider("user")
	assert.NotNil(t, provider)

	_, nok := provider.GetSecret("test")
	assert.False(t, nok)

	err := provider.Configure(nil)
	assert.Error(t, err)

	cfg := new(testConfig)
	cfg.VaultCfg = map[string]interface{}{
		"address": s.URL,
		"app":     "app",
	}
	err = provider.Configure(cfg.VaultCfg)
	assert.NoError(t, err)

	_, ok := provider.GetSecret("test")
	assert.True(t, ok)

}
