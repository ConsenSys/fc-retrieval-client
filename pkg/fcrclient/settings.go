package fcrclient

// Copyright (C) 2020 ConsenSys Software Inc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

const (
	defaultSettingsFileName = "fcrsettings.json"
	defaultSettingsFile = "/etc/client/" + defaultSettingsFileName
	
	defaultSettingsVerbose = true
	defaultMaxEstablishmentTTL = int64(100)
)

// FilecoinRetrievalClientSettings holds the library configuration
type FilecoinRetrievalClientSettings struct {
	MaxEstablishmentTTL int64 
	Verbose bool       // If true, then more logging is shown.
	NodeID *nodeid.NodeID
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

	// TODO handle node ID
	nodeID := nodeid.NewNodeID(big.NewInt(7))
	settings.NodeID = nodeID

	if settings.Verbose {
		log.Printf("Filecoin Retrieval Client settings: (%+v)\n", settings)
	}

	return settings, nil
}

// SetSettings allows the settings object to be created in memory
func SetSettings(set *FilecoinRetrievalClientSettings) {
	settings = set
}
