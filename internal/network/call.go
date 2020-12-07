package network



import (
	"github.com/bitly/go-simplejson"
	"log"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
    //"fmt"
)

const (
	apiURL string = "http://gateway:80/client/establishment" 
)

// GatewayCall calls the Gateway's REST API
func GatewayCall(method string, args map[string]interface{}) (*simplejson.Json) {
	args["protocol_version"] = "1"
	args["protocol_supported"] = "1"
	args["message_type"] = method
	mJSON, _ := json.Marshal(args)
	log.Printf("JSON sent: %s", mJSON)
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", apiURL, contentReader)
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
