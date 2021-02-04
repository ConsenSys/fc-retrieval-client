package gateway

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/request"
	"github.com/ConsenSys/fc-retrieval-register/pkg/register"
)

// GetRegisteredGateways registered Gateway list
func GetRegisteredGateways() ([]register.GatewayRegister, error) {
	url := "http://localhost:8080/registers/gateway"
	gateways := []register.GatewayRegister{}
	err := request.GetJSON(url, &gateways)
	if err != nil {
		log.Error("%+v", err)
		return gateways, err
	}
	if len(gateways) == 0 {
		log.Warn("No gateways found")
	}
	return gateways, nil
}
