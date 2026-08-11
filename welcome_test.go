// NOTE: This test file was written by AI (Claude).
package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceBuilderAddDNS(t *testing.T) {
	services := NewServiceBuilder().AddDNS(true, 53, nil).Build()

	require.Len(t, services, 1)
	assert.Equal(t, "dns", services[0].Name)
	assert.True(t, services[0].Enabled)
	assert.Equal(t, "127.0.0.1:53", services[0].Address)
	assert.NoError(t, services[0].Error)
}

func TestServiceBuilderAddDNSWithError(t *testing.T) {
	bindErr := errors.New("address already in use")

	services := NewServiceBuilder().AddDNS(true, 53, bindErr).Build()

	require.Len(t, services, 1)
	assert.Equal(t, bindErr, services[0].Error)
}

func TestServiceBuilderDisabled(t *testing.T) {
	services := NewServiceBuilder().AddDNS(false, 5353, nil).Build()

	require.Len(t, services, 1)
	assert.False(t, services[0].Enabled)
	assert.Equal(t, "127.0.0.1:5353", services[0].Address)
}

func TestServiceBuilderChaining(t *testing.T) {
	builder := NewServiceBuilder()

	// AddDNS returns the builder so future services can chain
	assert.Same(t, builder, builder.AddDNS(true, 53, nil))
}
