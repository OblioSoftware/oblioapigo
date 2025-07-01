package main

import (
	"fmt"
	"os"
	"time"

	oblioapi "github.com/obliosoftware/oblioapigo"
)

func main() {
	api := oblioapi.Api{
		ClientID:     os.Getenv("OBLIO_CLIENT_ID"),
		ClientSecret: os.Getenv("OBLIO_CLIENT_SECRET"),
	}
	response, err := api.CreateDoc("invoice", oblioapi.Doc{
		Cif: "1111",
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
				VatPercentage:            19,
				VatIncluded:              true,
				Quantity:                 1,
				ProductType:              "Serviciu",
				Management:               "CV",
				Save:                     false,
			},
		},
	})
	fmt.Printf("response %+v %+v\n", response, err)
}
