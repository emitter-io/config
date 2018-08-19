// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package vault

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/acme/autocert"
)

func TestCertCache(t *testing.T) {
	s := httptest.NewServer(&testVaultHandler{})
	defer s.Close()

	provider := NewProvider("user")
	assert.NotNil(t, provider)

	err := provider.Configure(map[string]interface{}{
		"address": s.URL,
		"app":     "test",
	})
	assert.NoError(t, err)

	cache := provider.GetCache()
	assert.NotNil(t, cache)

	{
		_, err := cache.Get(context.Background(), "aaa")
		assert.Error(t, err, autocert.ErrCacheMiss)
	}

	{
		b, err := cache.Get(context.Background(), "abc")
		assert.NoError(t, err)
		assert.Equal(t, "secret certificate", string(b))
	}

}
