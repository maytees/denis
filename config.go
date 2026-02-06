package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kirsle/configdir"
	"github.com/pelletier/go-toml/v2"
)

// A record under records.toml
type Record struct {
	Name  string `toml:"name"`
	Type  string `toml:"type"`
	Value string `toml:"value"`
	TTL   uint32 `toml:"ttl"`
}

// Holds many records ^^
type DNSRecords struct {
	Records []Record `toml:"records"`
}

// Holds DNS config of denis (seperate from proxy, api & web server)
type DNSConfig struct {
	Enabled  bool   `toml:"enabled"`
	Port     int    `toml:"port"`
	Upstream string `toml:"upstream"` // fallback
}

// Holds dns, proxy, api & web server
type DenisConfig struct {
	DNS DNSConfig `toml:"dns"`
}

// Gets config directory (not file)
func getConfigPath() string {
	// Checks current dir (for local dev)
	if _, err := os.Stat("./config"); err == nil {
		return "./config"
	}

	configPath := configdir.LocalConfig("denis")
	configdir.MakePath(configPath)
	return configPath
}

func loadDNSRecords() DNSRecords {
	configPath := getConfigPath()
	recordsFile := filepath.Join(configPath, "records.toml")

	if _, err := os.Stat(recordsFile); os.IsNotExist(err) {
		defaultRecords := DNSRecords{
			Records: []Record{
				{
					Name:  "example.local",
					Type:  "A",
					Value: "127.0.0.1",
					TTL:   300,
				},
			},
		}

		data, _ := toml.Marshal(defaultRecords)
		os.WriteFile(recordsFile, data, 0777)

		// No need to read
		return defaultRecords
	}

	// Read & unmarshal (deserialize)
	data, _ := os.ReadFile(recordsFile)
	var cfg DNSRecords
	toml.Unmarshal(data, &cfg)
	return cfg
}

func loadDNSConfig() DenisConfig {
	configPath := getConfigPath()
	configFile := filepath.Join(configPath, "config.toml")

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := DenisConfig{
			DNS: DNSConfig{
				Enabled:  true,
				Port:     53,
				Upstream: "8.8.8.8:53",
			},
		}

		data, _ := toml.Marshal(defaultConfig)
		os.WriteFile(configFile, data, 0777)

		// No need to read
		return defaultConfig
	}

	// Read & unmarshal (deserialize)
	data, _ := os.ReadFile(configFile)
	var cfg DenisConfig
	toml.Unmarshal(data, &cfg)
	return cfg
}

// TODO: Slow linear search, optimize later
func findRecordByName(records []Record, target string) (Record, bool) {
	for _, record := range records {
		// Records are always stored lowercase
		if record.Name == strings.ToLower(target) {
			return record, true
		}
	}

	return Record{}, false
}
