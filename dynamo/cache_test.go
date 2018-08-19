// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package dynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/acme/autocert"
)

func TestCertCache(t *testing.T) {
	ddb := new(mockDynamoDB)
	ddb.cache = make(map[string]map[string]*dynamodb.AttributeValue)
	provider := NewProvider()
	assert.NotNil(t, provider)

	err := provider.Configure(map[string]interface{}{
		"region": "ap-southeast-1",
		"table":  "test",
	})
	assert.NoError(t, err)

	provider.client = &client{
		dynamo:    ddb,
		table:     "test",
		keyColumn: "k",
		valColumn: "v",
	}

	cache, _ := provider.GetCache()
	assert.NotNil(t, cache)

	{
		err := cache.Put(context.Background(), "abc", []byte("secret certificate"))
		assert.NoError(t, err)
	}

	{
		b, err := cache.Get(context.Background(), "abc")
		assert.NoError(t, err)
		assert.Equal(t, "secret certificate", string(b))
	}

	{
		err := cache.Delete(context.Background(), "abc")
		assert.NoError(t, err)
	}

	{
		v, err := cache.Get(context.Background(), "abc")
		assert.Empty(t, v)
		assert.Error(t, err, autocert.ErrCacheMiss)
	}

}
