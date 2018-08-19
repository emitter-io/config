// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
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

	{
		err := provider.client.Put("mySecret", "hi there")
		assert.NoError(t, err)
	}

	{
		v, ok := provider.GetSecret("mySecret")
		assert.True(t, ok)
		assert.Equal(t, "hi there", v)
	}

	{
		v, ok := provider.GetSecret("someOtherSecret")
		assert.False(t, ok)
		assert.Empty(t, v)
	}

}
