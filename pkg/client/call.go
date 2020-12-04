package client



import (
	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"log"
    //"fmt"
)

const (
	apiUrl string = "http://gateway:80/client/establishment" 
)

func call(args ...string) (*simplejson.Json) {
	method := args[0]

	protocolVer := `protocol_version":"1"`
	protocolSupported := `"protocol_supported":"1"`
	messageType := `message_type":"` + method + `"` 
	postBodyPart1 := `{"jsonrpc":"2.0","` + 
		protocolVer + `,` + 
		protocolSupported + `,` + 
		messageType

	postBodyPart2 := ""
	for i:= 1; i < len(args); i++ {
		postBodyPart2 = postBodyPart2 + "," + args[i]
	}
	
	postBodyPart3 := `}`
	postBody := postBodyPart1 + postBodyPart2 + postBodyPart3
	log.Println("postBody: " + postBody)

	_, body, errs := gorequest.New().Post(apiUrl).
		Send(postBody).
		End()

	if errs != nil {
		panic(errs)
	}
	log.Println("response body: " + body)


	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		log.Fatalln(err)
	}

	return js
}
