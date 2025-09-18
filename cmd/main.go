package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	oblioapi "github.com/obliosoftware/oblioapigo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api := oblioapi.Api{
		TokenHandler: &oblioapi.TokenHandler{
			ClientID:     os.Getenv("OBLIO_CLIENT_ID"),
			ClientSecret: os.Getenv("OBLIO_CLIENT_SECRET"),
		},
	}

	response, err := api.Nomenclature("companies", map[string]string{})
	fmt.Printf("companies response %+v %+v\n", response, err)

	response, err = api.Nomenclature("clients", map[string]string{
		"cif":       os.Getenv("OBLIO_CIF"),
		"clientCif": "1111",
	})
	fmt.Printf("clients response %+v %+v\n", response, err)

	response, err = api.CreateDoc("invoice", oblioapi.Doc{
		Cif: os.Getenv("OBLIO_CIF"),
		Client: oblioapi.Client{
			Cif:          "",
			Name:         "Irina Fabiola",
			Rc:           "",
			Code:         "",
			Address:      "Progresul Bloc 32, Numarul 5",
			State:        "Brasov",
			City:         "Brasov",
			Country:      "Romania",
			Email:        "",
			Phone:        "",
			Contact:      "Irina Fabiola",
			Save:         true,
			Autocomplete: false,
		},
		SeriesName: "FCT",
		IssueDate:  time.Now(),
		DueDate:    time.Now(),
		Products: []oblioapi.Product{
			{
				Name:                     "Hemograma cu formula leucocitara, Hb,Ht,indici si reticulocite (Hemograma)",
				Code:                     "",
				Description:              "",
				Price:                    49.5,
				MeasuringUnit:            "buc",
				MeasuringUnitTranslation: "",
				Currency:                 "RON",
				VatPercentage:            21,
				VatIncluded:              true,
				Quantity:                 1,
				ProductType:              "Serviciu",
				Management:               "CV",
				Save:                     false,
			},
		},
	})
	fmt.Printf("doc response %+v %+v\n", response, err)
}
