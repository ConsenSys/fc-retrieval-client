package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)


// BuilderImpl holds the library configuration
type BuilderImpl struct {
	logLevel     string
	logTarget    string
	establishmentTTL int64
	clientID              *nodeid.NodeID

	blockchainPrivateKey	*fcrcrypto.KeyPair

	retrievalPrivateKey		*fcrcrypto.KeyPair
	retrievalPrivateKeyVer	*fcrcrypto.KeyVersion
}

// CreateSettings creates an object with the default settings.
func CreateSettings() (*BuilderImpl) {
	f := BuilderImpl{}
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
	f.establishmentTTL = defaultEstablishmentTTL
	return &f
}


// SetLogging sets the log level and target.
func (f *BuilderImpl) SetLogging(logLevel string, logTarget string) {
	f.logLevel = defaultLogLevel
	f.logTarget = defaultLogTarget
}

// SetEstablishmentTTL sets the time to live for the establishment message between client and gateway.
func (f *BuilderImpl) SetEstablishmentTTL(ttl int64) {
	f.establishmentTTL = ttl
}

// SetBlockchainPrivateKey sets the blockchain private key.
func (f *BuilderImpl) SetBlockchainPrivateKey(bcPkey *fcrcrypto.KeyPair) {
	f.blockchainPrivateKey = bcPkey
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f *BuilderImpl) SetRetrievalPrivateKey(rPkey *fcrcrypto.KeyPair, ver *fcrcrypto.KeyVersion) {
	f.retrievalPrivateKey = rPkey
	f.retrievalPrivateKeyVer = ver
}


// Build creates a settings object and initialises the logging system.
func (f *BuilderImpl) Build() (*ClientSettings){
	var err error

	logging.SetLogLevel(f.logLevel)
	logging.SetLogTarget(f.logTarget)

	g := ClientSettings{}
	g.establishmentTTL = f.establishmentTTL

	if f.blockchainPrivateKey == nil {
		logging.ErrorAndPanic("Settings: Blockchain Private Key not set")
	}
	g.blockchainPrivateKey = f.blockchainPrivateKey

	if f.clientID == nil {
		logging.Info("Settings: No Client ID set. Generating random client ID")
		// TODO replace once NewRandomNodeID becomes available.
		g.clientID , err = nodeid.NewRandomNodeID()
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random client ID: %s", err)
		}
	} else {
		g.clientID = f.clientID 
	}

	if (f.retrievalPrivateKey == nil) {
		pKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
		if (err != nil) {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s", err)
		}
		f.retrievalPrivateKey = pKey
		f.retrievalPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.retrievalPrivateKey = f.retrievalPrivateKey
		g.retrievalPrivateKeyVer = f.retrievalPrivateKeyVer
	}

	return &g
}




