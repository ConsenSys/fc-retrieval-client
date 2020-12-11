package gatewayapi



import (
	"github.com/bitly/go-simplejson"
	"log"
	"encoding/json"
	"bytes"
	"errors"
	"net"
	"net/http"
	"io/ioutil"
    //"fmt"
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
}

// NewGatewayAPIComms creates a connection with a gateway
func NewGatewayAPIComms(host string) (*Comms, error){
	// Create the constant array.
	if (clientAPIProtocolSupported == nil) {
		clientAPIProtocolSupported = make([]int, 1)
		clientAPIProtocolSupported[0] = clientAPIProtocolSupportedHi
	}

	// Check that the host name is valid
	err := validateHostName(host)
	if (err != nil) {
		log.Println(err)
		return nil, err
	}

	netComms := Comms{}
	netComms.apiURL = apiURLStart + host + apiURLEnd
	return &netComms, nil
}

// GatewayCall calls the Gateway's REST API
func (n *Comms) gatewayCall(method int32, args map[string]interface{}) (*simplejson.Json) {
	args["protocol_version"] = int32(1)
	args["protocol_supported"] = []int32{1}
	args["message_type"] = method
	mJSON, _ := json.Marshal(args)
	log.Printf("JSON sent: %s", mJSON)
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", n.apiURL, contentReader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errs := client.Do(req)
	if errs != nil {
		panic(errs)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	log.Printf("response body: %s", string(data))

	js, err := simplejson.NewJson(data)
	if err != nil {
		log.Fatalln(err)
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
	log.Printf("Resolved %s as %s\n", host, ra.String())
	return nil
}