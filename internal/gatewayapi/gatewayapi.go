package gatewayapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/ConsenSys/fc-retrieval-client/internal/contracts"
	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/bitly/go-simplejson"
)

const (
	apiURLStart string = "http://" 
	apiURLEnd string = "/v1" 
//	apiURLEnd string = "/client/establishment" 
)


const (
	clientAPIProtocolVersion = 1
	clientAPIProtocolSupportedHi = 1
)
// Can't have constant slices so create this at runtime.
// Order the API versions from most desirable to least desirable.
var clientAPIProtocolSupported []int



// Comms holds the communications specific data
type Comms struct {
	apiURL string
	gatewayPubKey *fcrcrypto.KeyPair
	gatewayPubKeyVer *fcrcrypto.KeyVersion
	settings *settings.ClientSettings
}

// NewGatewayAPIComms creates a connection with a gateway
func NewGatewayAPIComms(gatewayInfo *contracts.GatewayInformation, settings *settings.ClientSettings) (*Comms, error){
	host := gatewayInfo.Hostname
	port := gatewayInfo.HostPort

	// Create the constant array.
	if (clientAPIProtocolSupported == nil) {
		clientAPIProtocolSupported = make([]int, 1)
		clientAPIProtocolSupported[0] = clientAPIProtocolSupportedHi
	}

	// Check that the host name is valid
	err := validateHostName(host)
	if (err != nil) {
		logging.Error("Host name invalid: %s", err. Error())
		return nil, err
	}

	netComms := Comms{}
	netComms.apiURL = apiURLStart + host + ":" + strconv.Itoa(port) + apiURLEnd
	netComms.gatewayPubKey = gatewayInfo.GatewayRetrievalPublicKey
	netComms.gatewayPubKeyVer = gatewayInfo.GatewayRetrievalPublicKeyVersion
	netComms.settings = settings
	return &netComms, nil
}


// GatewayCall calls the Gateway's REST API
func (c *Comms) gatewayCall(msg interface{}) (*simplejson.Json) {

	// Create HTTP request.
	mJSON, _ := json.Marshal(msg)
	logging.Info("JSON sent: %s", string(mJSON))
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", c.apiURL, contentReader)
	req.Header.Set("Content-Type", "application/json")

	// Send request and receive response.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.ErrorAndPanic("Client - Gateway communications (%s): %s", c.apiURL, err)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	logging.Info("response body: %s", string(data))

	js, err := simplejson.NewJson(data)
	if err != nil {
		logging.ErrorAndPanic("Error decoding JSON: %s", err.Error())
	}

	return js
}


func validateHostName(host string) error {
	if len(host) == 0 {
		return errors.New("Error: Host name empty")
	} 

	ra, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return errors.New("Error: DNS Look-up failed for host: " + host)
	}
	logging.Info("Resolved %s as %s\n", host, ra.String())
	return nil
}


func (c *Comms) verifyMessage(signature string, wholeMessage interface{}) bool {
	keyVersion, err := fcrcrypto.ExtractKeyVersionFromMessage(signature)
	if err != nil {
		logging.Warn("Error decodign signature: %+v", err)
		return false
	}
	if keyVersion.NotEquals(c.gatewayPubKeyVer) {
		// TODO need to allow for multiple key versions, and fetch correct key
		// TODO based on version.
		logging.Error("Unknown Key Version used by gateway: %d", keyVersion.EncodeKeyVersion())
	}

	verified, err := fcrcrypto.VerifyMessage(c.gatewayPubKey, signature, wholeMessage)
	if err != nil {
		logging.Warn("Signature verification error: %+v", err)
		return false
	}
	return verified
}