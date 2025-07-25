package oblioapi

import (
	"encoding/json"
	"net/http"
	"time"
)

type AccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   uint   `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
	RequestTime uint   `json:"request_time,omitempty"`
}

func (v *AccessToken) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	ExpiresIn, _ := AnyToType[uint](item["expires_in"])
	RequestTime, _ := AnyToType[uint](item["request_time"])

	*v = AccessToken{
		AccessToken: StringFromInterface(item["access_token"]),
		ExpiresIn:   ExpiresIn,
		TokenType:   StringFromInterface(item["token_type"]),
		Scope:       StringFromInterface(item["scope"]),
		RequestTime: RequestTime,
	}

	return nil
}

func (v AccessToken) IsValid() bool {
	return (v.RequestTime + v.RequestTime) > uint(time.Now().UnixMilli()/1000)
}

type Response struct {
	Status        int    `json:"status,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Data          any    `json:"data,omitempty"`
}

type Doc struct {
	Cif                string    `json:"cif,omitempty"`
	IssueDate          time.Time `json:"issueDate,omitempty"`
	DueDate            time.Time `json:"dueDate,omitempty"`
	SeriesName         string    `json:"seriesName,omitempty"`
	Language           string    `json:"language,omitempty"`
	Precision          uint      `json:"precision,omitempty"`
	Currency           string    `json:"currency,omitempty"`
	IssuerName         string    `json:"issuerName,omitempty"`
	IssuerId           string    `json:"issuerId,omitempty"`
	NoticeNumber       string    `json:"noticeNumber,omitempty"`
	InternalNote       string    `json:"internalNote,omitempty"`
	DeputyName         string    `json:"deputyName,omitempty"`
	DeputyIdentityCard string    `json:"deputyIdentityCard,omitempty"`
	DeputyAuto         string    `json:"deputyAuto,omitempty"`
	SelesAgent         string    `json:"selesAgent,omitempty"`
	Mentions           string    `json:"mentions,omitempty"`
	WorkStation        string    `json:"workStation,omitempty"`
	Client             Client    `json:"client,omitempty"`
	Products           []Product `json:"products,omitempty"`
}

type Client struct {
	Cif          string `json:"cif,omitempty"`
	Name         string `json:"name,omitempty"`
	Rc           string `json:"rc,omitempty"`
	Code         string `json:"code,omitempty"`
	Address      string `json:"address,omitempty"`
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	Iban         string `json:"iban,omitempty"`
	Bank         string `json:"bank,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Contact      string `json:"contact,omitempty"`
	VatPayer     bool   `json:"vatPayer,omitempty"`
	Save         bool   `json:"save,omitempty"`
	Autocomplete bool   `json:"autocomplete,omitempty"`
}

type Product struct {
	ID                       uint    `json:"id,omitempty"`
	Name                     string  `json:"name,omitempty"`
	Description              string  `json:"description,omitempty"`
	Code                     string  `json:"code,omitempty"`
	MeasuringUnit            string  `json:"measuringUnit,omitempty"`
	MeasuringUnitTranslation string  `json:"MeasuringUnitTranslation,omitempty"`
	ProductType              string  `json:"productType,omitempty"`
	Price                    float64 `json:"price"`
	Quantity                 float64 `json:"quantity"`
	ExchangeRate             float64 `json:"exchangeRate"`
	VatName                  string  `json:"vatName,omitempty"`
	VatPercentage            float64 `json:"vatPercentage,omitempty"`
	VatIncluded              bool    `json:"vatIncluded"`
	Currency                 string  `json:"currency,omitempty"`
	Management               string  `json:"management,omitempty"`

	ProductID   uint    `json:"productId,omitempty"`
	Discount    float64 `json:"discount,omitempty"`
	DscountType string  `json:"discountType,omitempty"`
	Save        bool    `json:"save,omitempty"`
}

type TokenHandlerInterface interface {
	Get() (AccessToken, error)
}

type TokenHandler struct {
	ClientID     string
	ClientSecret string
	AccessToken  AccessToken
}

func (a *TokenHandler) Get() (AccessToken, error) {
	if a.AccessToken.IsValid() {
		return a.AccessToken, nil
	}
	response, err := Request(http.MethodPost, "/api/authorize/token", Payload{
		Type: "application/x-www-form-urlencoded",
		Data: map[string]any{
			"client_id":     a.ClientID,
			"client_secret": a.ClientSecret,
		},
	})
	if err != nil {
		return a.AccessToken, err
	}
	data, err := ReadResponse(response)
	if err != nil {
		return a.AccessToken, err
	}
	err = a.AccessToken.UnmarshalJSON(data)
	if err != nil {
		return a.AccessToken, err
	}

	return a.AccessToken, nil
}

type Api struct {
	TokenHandler TokenHandlerInterface
}

func (a Api) CreateDoc(t string, d Doc) (*Response, error) {
	token, err := a.TokenHandler.Get()
	if err != nil {
		return nil, err
	}

	response, err := Request(http.MethodPost, "/api/docs/"+t, Payload{
		AccessToken: &token,
		Type:        "application/json",
		Data:        d,
	})
	return a.getResponse(response, err)
}

func (a Api) Nomenclature(t string, filters map[string]string) (*Response, error) {
	token, err := a.TokenHandler.Get()
	if err != nil {
		return nil, err
	}

	response, err := Request(http.MethodGet, "/api/nomenclature/"+t+"?"+Urlencode(filters), Payload{
		AccessToken: &token,
		Type:        "application/json",
	})
	return a.getResponse(response, err)
}

func (a Api) getResponse(response *http.Response, err error) (*Response, error) {
	if err != nil {
		return nil, err
	}
	data, err := ReadResponse(response)
	if err != nil {
		return nil, err
	}
	res := Response{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
