package fcrclient

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-client/pkg/api/clientapi"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

// FilecoinRetrievalClient is an example implementation using the api,
// which holds information about the interaction of the Filecoin
// Retrieval Client with Filecoin Retrieval Gateways.
type FilecoinRetrievalClient struct {
	Settings ClientSettings

	// List of gateways this client can potentially use
	GatewaysToUse     map[string]register.GatewayRegister
	GatewaysToUseLock sync.RWMutex

	// List of gateway in use. A client may request a node be added to this list
	ActiveGateways     map[string]register.GatewayRegister
	ActiveGatewaysLock sync.RWMutex

	// PaymentMgr payment manager
	PaymentMgr *fcrpaymentmgr.FCRPaymentMgr
}

// NewFilecoinRetrievalClient initialise the Filecoin Retrieval Client library
func NewFilecoinRetrievalClient(settings ClientSettings) (*FilecoinRetrievalClient, error) {
	f := &FilecoinRetrievalClient{
		Settings:           settings,
		GatewaysToUse:      make(map[string]register.GatewayRegister),
		GatewaysToUseLock:  sync.RWMutex{},
		ActiveGateways:     make(map[string]register.GatewayRegister),
		ActiveGatewaysLock: sync.RWMutex{},
	}

	mgr, err := fcrpaymentmgr.NewFCRPaymentMgr(settings.walletPrivateKey, settings.lotusAP, settings.lotusAuthToken)
	if err != nil {
		logging.Error("Error initializing payment manager.")
		return nil, err
	}

	f.PaymentMgr = mgr

	return f, nil
}

// FindGateways find gateways located near to the specified location. Use AddGateways
// to use these gateways.
func (c *FilecoinRetrievalClient) FindGateways(location string, maxNumToLocate int) ([]*nodeid.NodeID, error) {
	// Determine gateways to use. For the moment, this is just "use all of them"
	// TODO: This will have to become, use gateways that this client has FIL registered with.
	gateways, err := register.GetRegisteredGateways(c.Settings.RegisterURL())
	if err != nil {
		logging.Error("Error in getting all registered gateways.")
		return nil, err
	}

	res := make([]*nodeid.NodeID, 0)
	count := 0
	for _, info := range gateways {
		if info.GetRegionCode() == location {
			nodeID, err := nodeid.NewNodeIDFromHexString(info.GetNodeID())
			if err != nil {
				logging.Error("Error in generating node id, skipping: %v", info.GetNodeID())
				continue
			}
			res = append(res, nodeID)
			count++
			if count >= maxNumToLocate {
				break
			}
		}
	}
	return res, nil
}

// AddGatewaysToUse adds one or more gateways to use.
func (c *FilecoinRetrievalClient) AddGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.GatewaysToUseLock.RLock()
		_, exist := c.GatewaysToUse[gwToAddID.ToString()]
		c.GatewaysToUseLock.RUnlock()
		if exist {
			continue
		}
		gateway, err := register.GetGatewayByID(c.Settings.RegisterURL(), gwToAddID)
		if err != nil {
			logging.Error("Error getting registered gateway %v: %v", gwToAddID, err.Error())
			continue
		}
		if !validateGatewayInfo(&gateway) {
			logging.Error("Register info not valid.")
			continue
		}
		// Success
		c.GatewaysToUseLock.Lock()
		c.GatewaysToUse[gwToAddID.ToString()] = gateway
		c.GatewaysToUseLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveGatewaysToUse removes one or more gateways from the list of Gateways to use.
// This also removes the gateway from gateways in active map.
func (c *FilecoinRetrievalClient) RemoveGatewaysToUse(gwNodeIDs []*nodeid.NodeID) int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.GatewaysToUse[gwToRemoveID.ToString()]
		if exist {
			delete(c.GatewaysToUse, gwToRemoveID.ToString())
			numRemoved++
			c.ActiveGatewaysLock.Lock()
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			c.ActiveGatewaysLock.Unlock()
		}
	}

	return numRemoved
}

// RemoveAllGatewaysToUse removes all gateways from the list of Gateways.
// This also cleared all gateways in active
func (c *FilecoinRetrievalClient) RemoveAllGatewaysToUse() int {
	c.GatewaysToUseLock.Lock()
	defer c.GatewaysToUseLock.Unlock()
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.GatewaysToUse)
	c.GatewaysToUse = make(map[string]register.GatewayRegister)
	c.ActiveGateways = make(map[string]register.GatewayRegister)

	return numRemoved
}

// GetGatewaysToUse returns the list of gateways to use.
func (c *FilecoinRetrievalClient) GetGatewaysToUse() []*nodeid.NodeID {
	c.GatewaysToUseLock.RLock()
	defer c.GatewaysToUseLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.GatewaysToUse {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}

// AddActiveGateways adds one or more gateways to active gateway map.
// Returns the number of gateways added.
func (c *FilecoinRetrievalClient) AddActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	numAdded := 0
	for _, gwToAddID := range gwNodeIDs {
		c.ActiveGatewaysLock.RLock()
		_, exist := c.ActiveGateways[gwToAddID.ToString()]
		c.ActiveGatewaysLock.RUnlock()
		if exist {
			continue
		}
		c.GatewaysToUseLock.RLock()
		info, exist := c.GatewaysToUse[strings.ToLower(gwToAddID.ToString())]
		c.GatewaysToUseLock.RUnlock()
		if !exist {
			logging.Error("Given node id: %v does not exist in gateways to use map, consider add the gateway first.", gwToAddID.ToString())
			continue
		}
		// Attempt an establishment
		challenge := make([]byte, 32)
		rand.Read(challenge)
		ttl := time.Now().Unix() + c.Settings.EstablishmentTTL()
		err := clientapi.RequestEstablishment(&info, challenge, c.Settings.ClientID(), ttl)
		if err != nil {
			logging.Error("Error in initial establishment: %v", err.Error())
			continue
		}
		// It is success
		c.ActiveGatewaysLock.Lock()
		c.ActiveGateways[gwToAddID.ToString()] = info
		c.ActiveGatewaysLock.Unlock()
		numAdded++
	}
	return numAdded
}

// RemoveActiveGateways removes one or more gateways from the list of Gateways in active.
func (c *FilecoinRetrievalClient) RemoveActiveGateways(gwNodeIDs []*nodeid.NodeID) int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := 0
	for _, gwToRemoveID := range gwNodeIDs {
		_, exist := c.ActiveGateways[gwToRemoveID.ToString()]
		if exist {
			delete(c.ActiveGateways, gwToRemoveID.ToString())
			numRemoved++
		}
	}

	return numRemoved
}

// RemoveAllActiveGateways removes all gateways from the list of Gateways in active.
func (c *FilecoinRetrievalClient) RemoveAllActiveGateways() int {
	c.ActiveGatewaysLock.Lock()
	defer c.ActiveGatewaysLock.Unlock()

	numRemoved := len(c.ActiveGateways)
	c.ActiveGateways = make(map[string]register.GatewayRegister)

	return numRemoved
}

// GetActiveGateways returns the list of gateways that are active.
func (c *FilecoinRetrievalClient) GetActiveGateways() []*nodeid.NodeID {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	res := make([]*nodeid.NodeID, 0)
	for key := range c.ActiveGateways {
		nodeID, err := nodeid.NewNodeIDFromHexString(key)
		if err != nil {
			logging.Error("Error in generating node id.")
			continue
		}
		res = append(res, nodeID)
	}

	return res
}

// FindOffersStandardDiscovery finds offer using standard discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersStandardDiscovery(contentID *cid.ContentID, gatewayID *nodeid.NodeID) ([]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return make([]cidoffer.SubCIDOffer, 0), errors.New("Given gatewayID is not in active nodes map")
	}
	// TODO need to do nonce management
	offers, err := clientapi.RequestStandardDiscover(&gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), "", "")
	if err != nil {
		logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.NodeID, err)
		return make([]cidoffer.SubCIDOffer, 0), errors.New("Error in requesting standard discovery")
	}
	// Verify the offer one by one
	for _, offer := range offers {
		// Get provider's pubkey
		providerInfo, err := register.GetProviderByID(c.Settings.RegisterURL(), offer.GetProviderID())
		if !validateProviderInfo(&providerInfo) {
			logging.Error("Provider register info not valid.")
			continue
		}
		if err != nil {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		pubKey, err := providerInfo.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		// Verify the offer sig
		if offer.Verify(pubKey) != nil {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		// Now Verify the merkle proof
		if offer.VerifyMerkleProof() != nil {
			logging.Error("Merkle proof verification failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")
	}
	return offers, nil
}

// FindOffersDHTDiscovery finds offer using dht discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersDHTDiscovery(contentID *cid.ContentID, gatewayID *nodeid.NodeID, numDHT int64) (map[string]*[]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	offersMap := make(map[string]*[]cidoffer.SubCIDOffer)

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return offersMap, errors.New("Given gatewayID is not in active nodes map")
	}
	// TODO need to do nonce management
	contacted, contactedResp, uncontactable, err := clientapi.RequestDHTDiscover(&gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), numDHT, false, "", "")
	if err != nil {
		logging.Warn("GatewayDHTDiscovery error. Gateway: %s, Error: %s", gw.NodeID, err)
		return offersMap, errors.New("Error in requesting dht discovery")
	}
	for i := 0; i < len(uncontactable); i++ {
		logging.Warn("Gateway: %v is uncontactable.", uncontactable[i].ToString())
	}

	for i := 0; i < len(contacted); i++ {
		id := contacted[i]
		resp := contactedResp[i]
		// Verify the sub response
		// Get gateway's pubkey
		gatewayInfo, err := register.GetGatewayByID(c.Settings.RegisterURL(), &id)
		if err != nil {
			logging.Error("Error in getting gateway info.")
			continue
		}
		if !validateGatewayInfo(&gatewayInfo) {
			logging.Error("Gateway register info not valid.")
			continue
		}
		pubKey, err := gatewayInfo.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		if resp.Verify(pubKey) != nil {
			logging.Error("Fail to verify sub response.")
			continue
		}
		_, _, _, offers, _, err := fcrmessages.DecodeGatewayDHTDiscoverResponse(&resp)
		if err != nil {
			logging.Error("Fail to decode response")
			continue
		}
		entry := make([]cidoffer.SubCIDOffer, 0)
		for _, offer := range offers {
			// Get provider's pubkey
			providerInfo, err := register.GetProviderByID(c.Settings.RegisterURL(), offer.GetProviderID())
			if err != nil {
				logging.Error("Error in getting provider info.")
				continue
			}
			if !validateProviderInfo(&providerInfo) {
				logging.Error("Provider register info not valid.")
				continue
			}
			pubKey, err := providerInfo.GetSigningKey()
			if err != nil {
				logging.Error("Fail to obtain public key.")
				continue
			}
			// Verify the offer sig
			if offer.Verify(pubKey) != nil {
				logging.Error("Offer signature fail to verify.")
				continue
			}
			// Now Verify the merkle proof
			if offer.VerifyMerkleProof() != nil {
				logging.Error("Merkle proof verification failed.")
				continue
			}
			// Offer pass verification
			entry = append(entry, offer)
			logging.Info("Offer pass every verification, added to result")
		}
		offersMap[id.ToString()] = &entry
	}

	return offersMap, nil
}

// FindDHTOfferAck finds offer ack for a cid, gateway pair
func (c *FilecoinRetrievalClient) FindDHTOfferAck(contentID *cid.ContentID, gatewayID *nodeid.NodeID, providerID *nodeid.NodeID) (bool, error) {
	provider, err := register.GetProviderByID(c.Settings.RegisterURL(), providerID)
	if err != nil {
		logging.Error("Error getting registered provider %v: %v", providerID, err.Error())
		return false, errors.New("Provider not found inside register")
	}
	if !validateProviderInfo(&provider) {
		logging.Error("Register info not valid.")
		return false, errors.New("Invalid register info")
	}

	found, request, ack, err := clientapi.RequestDHTOfferAck(&provider, contentID, gatewayID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	// We need to verify	1, request is indeed a request, ack is indeed an ack.
	//						2, request is signed by provider, ack is signed by gateway
	//						3, ack and request has same nonce
	// TODO: We will have to also include the hash of the request inside the ack, so the hash can be checked

	// Get provider pub key and gw pub key
	// Get gateway's pubkey
	gatewayInfo, err := register.GetGatewayByID(c.Settings.RegisterURL(), gatewayID)
	if err != nil {
		logging.Error("Error in getting gateway info.")
		return false, errors.New("Error in getting gateway info")
	}
	if !validateGatewayInfo(&gatewayInfo) {
		logging.Error("Gateway register info not valid.")
		return false, errors.New("Gateway register info not valid")
	}
	gwPubKey, err := gatewayInfo.GetSigningKey()
	if err != nil {
		logging.Error("Fail to obtain public key.")
		return false, errors.New("Fail to obtain public key")
	}
	// Get provider's pubkey
	pvdPubKey, err := provider.GetSigningKey()
	if err != nil {
		logging.Error("Fail to obtain public key.")
		return false, errors.New("Fail to obtain public key")
	}
	// Verify the request.
	if request.Verify(pvdPubKey) != nil {
		return false, errors.New("Error in verifying request")
	}
	// Verify the offer indeed contains the given cid
	_, _, offers, err := fcrmessages.DecodeProviderPublishDHTOfferRequest(request)
	if err != nil {
		return false, err
	}
	found = false
	for _, offer := range offers {
		if offer.GetCIDs()[0].ToString() == contentID.ToString() {
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("Initial request does not contain the given cid")
	}
	// Verify the ack
	if ack.Verify(gwPubKey) != nil {
		return false, errors.New("Error in verifying the ack")
	}
	_, signature, err := fcrmessages.DecodeProviderPublishDHTOfferResponse(ack)
	if err != nil {
		return false, err
	}
	// Verify ack against request
	ok, err := fcrcrypto.VerifyMessage(gwPubKey, signature, request)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// FindOffersStandardDiscoveryV2 finds offer using standard discovery from given gateways
func (c *FilecoinRetrievalClient) FindOffersStandardDiscoveryV2(contentID *cid.ContentID, gatewayID *nodeid.NodeID, maxOffers int) ([]cidoffer.SubCIDOffer, error) {
	c.ActiveGatewaysLock.RLock()
	defer c.ActiveGatewaysLock.RUnlock()

	gw, exists := c.ActiveGateways[gatewayID.ToString()]
	if !exists {
		return make([]cidoffer.SubCIDOffer, 0), errors.New("Given gatewayID is not in active nodes map")
	}

	// It will call the payment manager to pay for the initial request
	paychAddrs, voucher, topup, cidOffers, err := c.pay(gw)
	if err != nil {
		return cidOffers, err
	}

	// There isn't any balance in the payment channel, need to topup (create)
	if topup == true {
		// If topup failed, then there is not enough balance, return error.
		err = c.PaymentMgr.Topup(gw.NodeID, c.Settings.TopUpAmount())
		if err != nil {
			logging.Warn("Topup. Gateway: %s, Error: %s", gw.NodeID, err)
			return make([]cidoffer.SubCIDOffer, 0), errors.New("Error in payment manager topup - is not enough balance")
		}
		paychAddrs, voucher, topup, cidOffers, err = c.pay(gw)
		if err != nil {
			return cidOffers, err
		}
	}

	// It pays for the first request to get a list of offer digests.
	// TODO need to do nonce management
	offerDigests, err := clientapi.RequestStandardDiscoverV2(&gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), paychAddrs, voucher)
	if err != nil {
		logging.Warn("GatewayStdDiscovery error. Gateway: %s, Error: %s", gw.NodeID, err)
		return make([]cidoffer.SubCIDOffer, 0), errors.New("Error in requesting standard discovery")
	}

	// Depend on the input (maximum num of offers) It will send further requests to request offers from the same gateway.
	offers, err := clientapi.RequestStandardDiscoverOffer(&gw, contentID, rand.Int63(), time.Now().Unix()+c.Settings.EstablishmentTTL(), offerDigests, paychAddrs, voucher)

	var validOffers []cidoffer.SubCIDOffer
	// Verify the offer one by one
	for _, offer := range offers {
		// Get provider's pubkey
		providerInfo, err := register.GetProviderByID(c.Settings.RegisterURL(), offer.GetProviderID())
		if !validateProviderInfo(&providerInfo) {
			logging.Error("Provider register info not valid.")
			continue
		}
		if err != nil {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		pubKey, err := providerInfo.GetSigningKey()
		if err != nil {
			logging.Error("Fail to obtain public key.")
			continue
		}
		// Verify the offer sig
		if offer.Verify(pubKey) != nil {
			logging.Error("Offer signature fail to verify.")
			continue
		}
		// Now Verify the merkle proof
		if offer.VerifyMerkleProof() != nil {
			logging.Error("Merkle proof verification failed.")
			continue
		}
		// Offer pass verification
		logging.Info("Offer pass every verification, added to result")

		validOffers = append(validOffers, offer)
		if len(validOffers) == maxOffers {
			break
		}
	}

	return validOffers, nil
}

func (c *FilecoinRetrievalClient) pay(gw register.GatewayRegister) (string, string, bool, []cidoffer.SubCIDOffer, error) {
	paychAddrs, voucher, topup, err := c.PaymentMgr.Pay(gw.NodeID, 0, c.Settings.searchPrice)
	if err != nil {
		logging.Warn("Pay error. Gateway: %s, Error: %s", gw.NodeID, err)
		return "", "", false, []cidoffer.SubCIDOffer{}, errors.New("Error in payment manager pay")
	}
	return paychAddrs, voucher, topup, nil, nil
}
