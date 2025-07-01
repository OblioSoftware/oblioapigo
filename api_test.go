package oblioapi_test

import (
	"os"
	"testing"

	oblioapi "github.com/obliosoftware/oblioapigo"
)

func TestConversions(t *testing.T) {
	api := oblioapi.Api{
		ClientID:     os.Getenv("OBLIO_CLIENT_ID"),
		ClientSecret: os.Getenv("OBLIO_CLIENT_SECRET"),
	}
	_, err := api.GetAccessToken()
	if err != nil {
		panic(err)
	}
}
