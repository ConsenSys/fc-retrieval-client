package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)

// ClientSettings holds the library configuration
type ClientSettings struct {
	establishmentTTL int64
	registerURL      string
	clientID         *nodeid.NodeID

	blockchainPrivateKey *fcrcrypto.KeyPair

	retrievalPrivateKey    *fcrcrypto.KeyPair
	retrievalPrivateKeyVer *fcrcrypto.KeyVersion
}

// EstablishmentTTL returns the establishmentTTL
func (c ClientSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// RegisterURL returns the register URL
func (c ClientSettings) RegisterURL() string {
	return c.registerURL
}

// ClientID returns the ClientID
func (c ClientSettings) ClientID() *nodeid.NodeID {
	return c.clientID
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c ClientSettings) BlockchainPrivateKey() *fcrcrypto.KeyPair {
	return c.blockchainPrivateKey
}

// RetrievalPrivateKey returns the RetrievalPrivateKey
func (c ClientSettings) RetrievalPrivateKey() *fcrcrypto.KeyPair {
	return c.retrievalPrivateKey
}

// RetrievalPrivateKeyVer returns the RetrievalPrivateKeyVer
func (c ClientSettings) RetrievalPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.retrievalPrivateKeyVer
}
