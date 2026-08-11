// NOTE: This test file was written by AI (Claude).
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVariantV4(t *testing.T) {
	assert.Equal(t, "v4", GetVariant("127.0.0.1"))
	assert.Equal(t, "v4", GetVariant("8.8.8.8"))
}

func TestGetVariantV6(t *testing.T) {
	assert.Equal(t, "v6", GetVariant("::1"))
	assert.Equal(t, "v6", GetVariant("2001:4860:4860::8888"))
}

func TestGetVariantDomain(t *testing.T) {
	assert.Equal(t, "domain", GetVariant("dns.google"))
	assert.Equal(t, "domain", GetVariant("localhost"))

	// Quirk: an invalid IP octet falls through to "domain"
	assert.Equal(t, "domain", GetVariant("256.1.1.1"))
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		host    string
		port    string
		variant string
		wantErr bool
	}{
		{"v4 with port", "8.8.8.8:53", "8.8.8.8", "53", "v4", false},
		{"v4 without port", "8.8.8.8", "8.8.8.8", "", "v4", false},
		{"v6 bracketed with port", "[::1]:53", "::1", "53", "v6", false},
		{"v6 without port", "::1", "::1", "", "v6", false},
		{"domain with port", "dns.google:53", "dns.google", "53", "domain", false},
		{"domain without port", "dns.google", "dns.google", "", "domain", false},
		{"garbage with colons", "a:b:c", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, variant, err := ParseAddress(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.host, host)
			assert.Equal(t, tt.port, port)
			assert.Equal(t, tt.variant, variant)
		})
	}
}

func TestBoolToUint16(t *testing.T) {
	assert.Equal(t, uint16(1), BoolToUint16(true))
	assert.Equal(t, uint16(0), BoolToUint16(false))
}
