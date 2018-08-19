// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package dynamo

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type mockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	cache map[string]map[string]*dynamodb.AttributeValue
}

func (m *mockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	if v, ok := m.cache[input.Key["k"].GoString()]; ok {
		return &dynamodb.GetItemOutput{
			Item: v,
		}, nil
	}

	return nil, errors.New("no such key")
}

func (m *mockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.cache[input.Item["k"].GoString()] = input.Item
	return &dynamodb.PutItemOutput{}, nil
}

func (m *mockDynamoDB) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	delete(m.cache, input.Key["k"].GoString())
	return &dynamodb.DeleteItemOutput{}, nil
}

func TestClient(t *testing.T) {
	ddb := new(mockDynamoDB)
	ddb.cache = make(map[string]map[string]*dynamodb.AttributeValue)
	client := &client{
		dynamo:    ddb,
		table:     "test",
		keyColumn: "k",
		valColumn: "v",
	}

	{
		err := client.Put("secret/a", "CAT")
		assert.NoError(t, err)
	}

	{
		v, err := client.Get("secret/a")
		assert.NoError(t, err)
		assert.Equal(t, "CAT", v)
	}

	{
		err := client.Delete("secret/a")
		assert.NoError(t, err)
	}

	{
		_, err := client.Get("secret/a")
		assert.Error(t, err)
	}
}
