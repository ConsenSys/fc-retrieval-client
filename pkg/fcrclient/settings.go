package fcrclient

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

const (
	defaultSettingsFileName = "fcrsettings.json"
	defaultSettingsFile     = "/etc/client/" + defaultSettingsFileName

	defaultMaxEstablishmentTTL = int64(100)

	settingsDefaultLogLevel = "TRACE"
    settingsDefaultLogTarget = "STDOUT"
)

// FilecoinRetrievalClientSettings holds the library configuration
type FilecoinRetrievalClientSettings struct {
	MaxEstablishmentTTL int64
	NodeID              *nodeid.NodeID
	LogLevel        string `json:"loglevel"`        // Log Level: NONE, ERROR, WARN, INFO, TRACE
	LogTarget       string `json:"logtarget"`       // Log Level: STDOUT
}

var settings *FilecoinRetrievalClientSettings

// LoadSettings loads the app settings from the settings file.
func LoadSettings(settingsFile ...string) (*FilecoinRetrievalClientSettings, error) {
	configFile := defaultSettingsFile
	if len(settingsFile) == 1 {
		configFile = settingsFile[0]
	}
	settingsBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("Failed to read settings file: %s: %s", configFile, err.Error())
		return nil, err
	}

	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		log.Printf("Failed to read settings file: %s: %s", configFile, err.Error())
	}

	logging.SetLogLevel(settings.LogLevel)
	logging.SetLogTarget(settings.LogTarget)
	logging.Info("Filecoin Retrieval Client settings: (%+v)", settings)
	return settings, nil
}

// SetSettings allows the settings object to be created in memory
func SetSettings(set *FilecoinRetrievalClientSettings) {
	settings = set
}
