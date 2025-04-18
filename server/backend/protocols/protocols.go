package protocols

import (
	"log"
	"path"
	"plugin"

	"github.com/ByteTheCookies/backend/internal/models"
)

func LoadProtocol(protocolName string) (func(string, string, []string) ([]models.ResponseProtocol, error), error) {
	pathProtocol := path.Join(".", "protocols", protocolName+".so")
	plug, err := plugin.Open(pathProtocol)
	if err != nil {
		log.Fatal(err)
	}

	sym, err := plug.Lookup("Submit")
	if err != nil {
		log.Fatal(err)
	}

	submitFunc := sym.(func(string, string, []string) ([]models.ResponseProtocol, error))

	return submitFunc, nil
}
