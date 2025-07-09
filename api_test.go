package oblioapi_test

import (
	"os"
	"testing"

	oblioapi "github.com/obliosoftware/oblioapigo"
)

func TestConversions(t *testing.T) {
	api := oblioapi.Api{
		TokenHandler: &oblioapi.TokenHandler{
			ClientID:     os.Getenv("OBLIO_CLIENT_ID"),
			ClientSecret: os.Getenv("OBLIO_CLIENT_SECRET"),
		},
	}
	_, err := api.TokenHandler.Get()
	if err != nil {
		panic(err)
	}
}

func TestToken(t *testing.T) {
	token := oblioapi.AccessToken{
		AccessToken: "57a17d5650f0633428ddc846979ef34c82c7dcdc",
		TokenType:   "Bearer",
		Scope:       "",
		ExpiresIn:   3600,
		RequestTime: 1752049737,
	}

	if !token.IsValid() {
		panic("Token not valid")
	}
}
