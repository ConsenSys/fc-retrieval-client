package gatewayapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/ConsenSys/fc-retrieval-client/internal/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
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
	settings *settings.ClientSettings
}

// NewGatewayAPIComms creates a connection with a gateway
func NewGatewayAPIComms(host string, settings *settings.ClientSettings) (*Comms, error){
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
	netComms.apiURL = apiURLStart + host + apiURLEnd
	netComms.settings = settings
	return &netComms, nil
}

// GatewayCall calls the Gateway's REST API
func (c *Comms) gatewayCall(method int32, args map[string]interface{}) (*simplejson.Json) {
	// Add in common fields to map.
	args["protocol_version"] = int32(1)
	args["protocol_supported"] = []int32{1}
	args["message_type"] = method
	args["node_id"] = c.settings.ClientID().ToString()

	// Sign fields.
	sig, err := fcrcrypto.SignMessage(c.settings.RetrievalPrivateKey(), c.settings.RetrievalPrivateKeyVer(), args)
	if err != nil {
		logging.ErrorAndPanic("Issue signing message: %s", err)
		panic(err)
	}
	args["signature"] = sig

	// Create HTTP request.
	mJSON, _ := json.Marshal(args)
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