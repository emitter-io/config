// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package config

import (
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
