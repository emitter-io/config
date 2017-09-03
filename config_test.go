package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testConfig struct {
	Name     string       `json:"name"`
	VaultCfg *VaultConfig `json:"vault,omitempty"`
}

func (c *testConfig) Vault() *VaultConfig {
	return c.VaultCfg
}

type secretStoreMock struct {
	mock.Mock
}

func (m *secretStoreMock) GetSecret(secretName string) (string, bool) {
	mockArgs := m.Called(secretName)
	v := mockArgs.Get(0).(string)
	return v, v != ""
}

func (m *secretStoreMock) Configure(c Config) error {
	return nil
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

	declassify(c, "emitter", m)

	assert.EqualValues(t, expected, c)

	c.Vault().Application = "abc"
	assert.True(t, c.Vault() != nil)
	assert.True(t, c.Vault().Application == "abc")
}
