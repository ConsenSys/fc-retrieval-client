package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

import (
	"crypto/ecdsa"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
)




// ClientSettings holds the library configuration
type ClientSettings struct {
	establishmentTTL int64
	clientID              *nodeid.NodeID

	blockchainPrivateKey	*ecdsa.PrivateKey 
	blockchainPrivateKeyAlg	*fcrcrypto.SigAlg

	retrievalPrivateKey		*ecdsa.PrivateKey
	retrievalPrivateKeyVer	*fcrcrypto.KeyVersion
	retrievalPrivateKeyAlg	*fcrcrypto.SigAlg
}


// EstablishmentTTL returns the establishmentTTL
func (c ClientSettings) EstablishmentTTL() int64 {
	return c.establishmentTTL
}

// ClientID returns the ClientID
func (c ClientSettings) ClientID() *nodeid.NodeID {
	return c.clientID
}

// BlockchainPrivateKey returns the BlockchainPrivateKey
func (c ClientSettings) BlockchainPrivateKey() *ecdsa.PrivateKey {
	return c.blockchainPrivateKey
}

// BlockchainPrivateKeyAlg returns the BlockchainPrivateKeyAlg
func (c ClientSettings) BlockchainPrivateKeyAlg() *fcrcrypto.SigAlg {
	return c.blockchainPrivateKeyAlg
}

// RetrievalPrivateKey returns the RetrievalPrivateKey
func (c ClientSettings) RetrievalPrivateKey() *ecdsa.PrivateKey {
	return c.retrievalPrivateKey
}

// RetrievalPrivateKeyVer returns the RetrievalPrivateKeyVer
func (c ClientSettings) RetrievalPrivateKeyVer() *fcrcrypto.KeyVersion {
	return c.retrievalPrivateKeyVer
}

// RetrievalPrivateKeyAlg returns the RetrievalPrivateKeyAlg
func (c ClientSettings) RetrievalPrivateKeyAlg() *fcrcrypto.SigAlg {
	return c.retrievalPrivateKeyAlg
}


