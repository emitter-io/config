// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	f, err := httpFile("http://google.com")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	b, err := ioutil.ReadFile(f.Name())
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	os.Remove(f.Name())
}
