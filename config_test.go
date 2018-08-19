// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package config

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testConfig struct {
	Name       string          `json:"name"`
	VaultCfg   *VaultConfig    `json:"vault,omitempty"`
	Provider   *ProviderConfig `json:"provider,omitempty"`
	unexported *ProviderConfig
}

type VaultConfig struct {
	Address     string `json:"address"` // The vault address to use.
	Application string `json:"app"`     // The vault application ID to use.
}

type secretStoreMock struct {
	mock.Mock
}

func (m *secretStoreMock) GetSecret(secretName string) (string, bool) {
	mockArgs := m.Called(secretName)
	v := mockArgs.Get(0).(string)
	return v, v != ""
}

func (m *secretStoreMock) Name() string {
	return "vault"
}

func (m *secretStoreMock) Configure(config map[string]interface{}) error {
	if config["address"] == nil || config["address"] == "" {
		return errors.New("address was not configured")
	}

	return nil
}

func TestReadOrCreate(t *testing.T) {
	m := new(secretStoreMock)
	m.On("GetSecret", "emitter/listen").Return(":999")
	m.On("GetSecret", "emitter/vault/address").Return("hello")
	m.On("GetSecret", mock.Anything).Return("")

	defaultCfg := new(testConfig)
	defaultCfg.Name = "test"
	defaultCfg.VaultCfg = &VaultConfig{
		Address: "test",
	}

	const f = "test.cfg"
	cfg, err := ReadOrCreate("test", f, func() Config { return defaultCfg }, m)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	os.Remove(f)
}

func Test_write(t *testing.T) {
	c := &testConfig{
		Name: ":80",
	}

	o := bytes.NewBuffer([]byte{})
	write(c, o)
	assert.Equal(t, "{\n\t\"name\": \":80\"\n}", string(o.Bytes()))
}

func Test_declassify(t *testing.T) {
	c := new(testConfig)
	c.Name = "test"
	c.VaultCfg = new(VaultConfig)

	m := new(secretStoreMock)
	m.On("GetSecret", "emitter/listen").Return(":999")
	m.On("GetSecret", "emitter/vault/address").Return("hello")
	m.On("GetSecret", mock.Anything).Return("")

	expected := new(testConfig)
	expected.Name = "test"
	expected.VaultCfg = new(VaultConfig)
	expected.VaultCfg.Address = "hello"
	expected.Provider = new(ProviderConfig)

	declassify(c, "emitter", m)

	assert.EqualValues(t, expected, c)
}

func Test_declassify_Map(t *testing.T) {
	c := new(testConfig)
	c.Name = "test"

	subconf := `{
		"type": "service_account",
		"project_id": "emitter-io",
		"private_key_id": "ABECESGSEGSEGSEG",
		"client_email": "EG#EFSD16546@appspot.gserviceaccount.com",
		"client_id": "65498913215443546542211557",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://accounts.google.com/o/oauth2/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/emitter-io%40appspot.gserviceaccount.com"}`

	m := new(secretStoreMock)
	m.On("GetSecret", "emitter/provider/config").Return(subconf)
	m.On("GetSecret", mock.Anything).Return("")

	declassify(c, "emitter", m)

	assert.Equal(t, "65498913215443546542211557", c.Provider.Config["client_id"])
}
