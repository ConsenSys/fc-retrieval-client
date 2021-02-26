package fcrclient

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"

	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
)

// SettingsBuilder holds the library configuration
type SettingsBuilder interface {
	// SetLogging sets the log level and target.
	SetLogging(logLevel string, logTarget string)

	// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
	SetEstablishmentTTL(ttl int64)

	// SetBlockchainPrivateKey sets the blockchain private key.
	SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair)

	// SetRetrievalPrivateKey sets the retrieval private key.
	SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion)

	// Build creates a settings object and initialises the logging system.
	Build() (*Settings)
}


// Settings holds the library configuration
type Settings interface {
	EstablishmentTTL() 		  int64
	ClientID() 				  *nodeid.NodeID

	BlockchainPrivateKey()    *fcrcrypto.KeyPair 

	RetrievalPrivateKey()	  *fcrcrypto.KeyPair
	RetrievalPrivateKeyVer()  *fcrcrypto.KeyVersion
}



// CreateSettings loads up default settings
func CreateSettings() (SettingsBuilder) {
	f := newBuilderImpl()
	builder := SettingsBuilder(f)
	return builder
}



type settingsBuilderImpl struct {
	impl *settings.BuilderImpl
}

func newBuilderImpl() settingsBuilderImpl {
  	return settingsBuilderImpl{settings.CreateSettings()}
}


// SetLogging sets the log level and target.
func (f settingsBuilderImpl) SetLogging(logLevel string, logTarget string) {
	f.impl.SetLogging(logLevel, logTarget)
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
func (f settingsBuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.impl.SetEstablishmentTTL(ttl)
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f settingsBuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.impl.SetBlockchainPrivateKey(bcPkey)
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f settingsBuilderImpl) SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.impl.SetRetrievalPrivateKey(rPkey, ver)
}

// Build generates the settings.
func (f settingsBuilderImpl) Build() *Settings {
	clientSettings := f.impl.Build()
	set := Settings(clientSettings)
	return &set
}

