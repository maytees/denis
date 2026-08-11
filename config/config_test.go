// NOTE: This test file was written by AI (Claude).
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindRecordByNameExact(t *testing.T) {
	records := []Record{
		{Name: "my.google", Type: "A", Value: "142.251.16.138", TTL: 0},
		{Name: "localhost", Type: "A", Value: "127.0.0.1", TTL: 300},
	}

	record, ok := FindRecordByName(records, "localhost")

	require.True(t, ok)
	assert.Equal(t, "127.0.0.1", record.Value)
}

func TestFindRecordByNameUppercaseQuery(t *testing.T) {
	records := []Record{
		{Name: "my.google", Type: "A", Value: "142.251.16.138", TTL: 0},
	}

	// DNS names are case-insensitive (RFC 1035 §2.3.3)
	record, ok := FindRecordByName(records, "MY.GOOGLE")

	require.True(t, ok)
	assert.Equal(t, "142.251.16.138", record.Value)
}

// GAP: records are assumed to be stored lowercase, but nothing normalizes
// them on load. A user writing 'My.Google' in records.toml gets a record
// that can never match any query.
func TestFindRecordByNameMixedCaseStored(t *testing.T) {
	records := []Record{
		{Name: "My.Google", Type: "A", Value: "142.251.16.138", TTL: 0},
	}

	_, ok := FindRecordByName(records, "my.google")

	assert.True(t, ok, "a record stored with mixed case should still be findable")
}

func TestFindRecordByNameMissing(t *testing.T) {
	records := []Record{
		{Name: "localhost", Type: "A", Value: "127.0.0.1", TTL: 300},
	}

	_, ok := FindRecordByName(records, "doesnotexist")

	assert.False(t, ok)
}

func TestFindRecordByNameEmptySlice(t *testing.T) {
	_, ok := FindRecordByName([]Record{}, "localhost")

	assert.False(t, ok)
}

// Sets up a temp working dir containing a ./config folder so
// GetConfigPath resolves there instead of the real config dir.
func chdirWithConfigDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	confDir := filepath.Join(tmp, "config")
	require.NoError(t, os.Mkdir(confDir, 0755))
	t.Chdir(tmp)
	return confDir
}

func TestGetConfigPathPrefersLocalDir(t *testing.T) {
	chdirWithConfigDir(t)

	assert.Equal(t, "./config", GetConfigPath())
}

func TestLoadDNSRecordsFromFile(t *testing.T) {
	confDir := chdirWithConfigDir(t)

	toml := `
[[records]]
name = 'notes'
type = 'A'
value = '127.0.0.1'
ttl = 300

[[records]]
name = 'my.google'
type = 'A'
value = '142.251.16.138'
ttl = 0
`
	require.NoError(t, os.WriteFile(filepath.Join(confDir, "records.toml"), []byte(toml), 0644))

	records := LoadDNSRecords()

	require.Len(t, records, 2)
	assert.Equal(t, "notes", records[0].Name)
	assert.Equal(t, "127.0.0.1", records[0].Value)
	assert.Equal(t, uint32(300), records[0].TTL)
	assert.Equal(t, "my.google", records[1].Name)
}

func TestLoadDNSRecordsCreatesDefaultWhenMissing(t *testing.T) {
	confDir := chdirWithConfigDir(t)

	records := LoadDNSRecords()

	require.Len(t, records, 1)
	assert.Equal(t, "example.local", records[0].Name)

	// The default file should now exist on disk
	_, err := os.Stat(filepath.Join(confDir, "records.toml"))
	assert.NoError(t, err)
}

// Documents current behavior: unmarshal errors are silently swallowed,
// so a malformed records.toml yields zero records with no warning.
func TestLoadDNSRecordsMalformedTOML(t *testing.T) {
	confDir := chdirWithConfigDir(t)
	require.NoError(t, os.WriteFile(filepath.Join(confDir, "records.toml"), []byte("not [valid toml"), 0644))

	records := LoadDNSRecords()

	assert.Empty(t, records)
}

func TestLoadDNSConfigFromFile(t *testing.T) {
	confDir := chdirWithConfigDir(t)

	toml := `
[dns]
enabled = false
port = 5353
upstream = '1.1.1.1'
`
	require.NoError(t, os.WriteFile(filepath.Join(confDir, "config.toml"), []byte(toml), 0644))

	cfg, path := LoadDNSConfig()

	assert.Equal(t, "./config", path)
	assert.False(t, cfg.DNS.Enabled)
	assert.Equal(t, 5353, cfg.DNS.Port)
	assert.Equal(t, "1.1.1.1", cfg.DNS.Upstream)
}

func TestLoadDNSConfigCreatesDefaultWhenMissing(t *testing.T) {
	confDir := chdirWithConfigDir(t)

	cfg, _ := LoadDNSConfig()

	assert.True(t, cfg.DNS.Enabled)
	assert.Equal(t, 53, cfg.DNS.Port)
	assert.Equal(t, "8.8.8.8:53", cfg.DNS.Upstream)

	_, err := os.Stat(filepath.Join(confDir, "config.toml"))
	assert.NoError(t, err)
}
