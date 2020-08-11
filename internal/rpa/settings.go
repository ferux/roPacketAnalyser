package rpa

import (
	"encoding/json"
	"fmt"
	"os"
)

// Settings are app-wide settings.
type Settings struct {
	Mode             AppMode  `json:"mode"`
	Addr             string   `json:"addr"`
	SupportedPackets []string `json:"supported_packets"`
	InclusiveFilter  bool     `json:"inclusive_filter"`
}

type AppMode string

const (
	AppModeDialer AppMode = "dial"
	AppModeServer AppMode = "server"
)

// ConfigFromFile reads config from file.
func ConfigFromFile(filepath string) (s Settings, err error) {
	var f *os.File
	f, err = os.Open(filepath)
	if err != nil {
		return Settings{}, fmt.Errorf("opening file: %w", err)
	}

	defer func() {
		_ = f.Close()
	}()

	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return Settings{}, fmt.Errorf("decoding into file: %w", err)
	}

	return s, nil
}
