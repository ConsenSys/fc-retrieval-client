package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"crypto/ecdsa"

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

	blockchainPrivateKey	*ecdsa.PrivateKey 
	blockchainPrivateKeyAlg	*fcrcrypto.SigAlg

	retrievalPrivateKey		*ecdsa.PrivateKey
	retrievalPrivateKeyVer	*fcrcrypto.KeyVersion
	retrievalPrivateKeyAlg	*fcrcrypto.SigAlg
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
func (f *BuilderImpl) SetBlockchainPrivateKey(bcPkey *ecdsa.PrivateKey, alg *fcrcrypto.SigAlg) {
	f.blockchainPrivateKey = bcPkey
	f.blockchainPrivateKeyAlg = alg
}

// SetRetrievalPrivateKey sets the retrieval private key.
func (f *BuilderImpl) SetRetrievalPrivateKey(rPkey *ecdsa.PrivateKey, alg *fcrcrypto.SigAlg, ver *fcrcrypto.KeyVersion) {
	f.retrievalPrivateKey = rPkey
	f.retrievalPrivateKeyAlg = alg
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
	g.blockchainPrivateKeyAlg = f.blockchainPrivateKeyAlg

	if f.clientID == nil {
		logging.Info("Settings: No Client ID set. Generating random client ID")
		// TODO replace once NewRandomNodeID becomes available.
//		g.clientID = nodeid.NewRandomNodeID()
		g.clientID, err = nodeid.NewNodeIDFromString("12345678")
		if err != nil {
			logging.ErrorAndPanic("Settings: Error while generating random client ID: %s", err)
		}
	} else {
		g.clientID = f.clientID 
	}

	if (f.retrievalPrivateKey == nil) {
		pKey, err := fcrcrypto.GenKeyPair()
		if (err != nil) {
			logging.ErrorAndPanic("Settings: Error while generating random retrieval key pair: %s", err)
		}
		f.retrievalPrivateKey = pKey
		f.retrievalPrivateKeyAlg = &fcrcrypto.SigAlg{Alg: fcrcrypto.SigAlgEcdsaP256Sha512_256}
		f.retrievalPrivateKeyVer = fcrcrypto.DecodeKeyVersion(1)
	} else {
		g.retrievalPrivateKey = f.retrievalPrivateKey
		g.retrievalPrivateKeyAlg = f.retrievalPrivateKeyAlg
		g.retrievalPrivateKeyVer = f.retrievalPrivateKeyVer
	}

	return &g
}




