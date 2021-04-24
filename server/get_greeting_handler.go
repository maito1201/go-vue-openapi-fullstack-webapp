package server

import (
	"log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/maito1201/go-vue-openapi-fullstack-webapp/server/gen/restapi/factory"
)

func GetGreeting(p factory.GetGreetingParams) middleware.Responder {
	payload := "hello go"
	if p.Name != nil {
		payload = *p.Name
	}
	log.Printf("GetGreeting is called, return %s\n", payload)
	return factory.NewGetGreetingOK().WithPayload(payload)
}
