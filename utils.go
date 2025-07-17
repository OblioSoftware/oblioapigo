package oblioapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://www.oblio.eu"
)

type Payload struct {
	Type        string
	Data        any
	AccessToken *AccessToken
}

func (v *Payload) Get() io.Reader {
	var p []byte
	switch v.Type {
	case "application/x-www-form-urlencoded":
		params := url.Values{}
		for i, v := range v.Data.(map[string]any) {
			params.Add(i, v.(string))
		}
		p = []byte(params.Encode())
	default:
		p, _ = json.Marshal(v.Data)
	}
	return bytes.NewReader(p)
}

func Request(method string, path string, payload Payload) (*http.Response, error) {
	req, err := http.NewRequest(method, baseURL+path, payload.Get())
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s\n", err)
	}
	req.Header.Add("Content-Type", payload.Type)
	if payload.AccessToken != nil {
		req.Header.Add("Authorization", payload.AccessToken.TokenType+" "+payload.AccessToken.AccessToken)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s\n", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		data, err := ReadResponse(res)
		if err != nil {
			return nil, err
		}
		var response map[string]any
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, err
		}
		message, ok := response["statusMessage"].(string)
		if !ok {
			message = "Http error"
		}
		return res, fmt.Errorf("%s", message)
	}

	return res, nil
}

func ReadResponse(res *http.Response) ([]byte, error) {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s\n", err)
	}
	return resBody, nil
}

func AnyToType[t uint | int | uint8 | int8 | int64 | float32 | float64](data ...interface{}) (t, error) {
	defaultVal := 0
	dataLen := len(data)
	if dataLen > 1 {
		defaultVal = data[1].(int)
	}
	switch data[0].(type) {
	case nil:
		return t(defaultVal), nil
	case string:
		val, err := strconv.ParseFloat(data[0].(string), 64)
		if err != nil && dataLen > 1 {
			return t(defaultVal), nil
		}
		return t(val), err
	case int8:
		return t(data[0].(int8)), nil
	case int:
		return t(data[0].(int)), nil
	case uint:
		return t(data[0].(uint)), nil
	case float64:
		return t(data[0].(float64)), nil
	default:
		return data[0].(t), nil
	}
}

func DateFromInterface(date interface{}) time.Time {
	dateString := ""
	if date != nil {
		dateString = date.(string)
	}
	if dateString == "" {
		dateString = "0001-01-01"
	}
	t, _ := time.Parse("2006-01-02", dateString[0:10])
	return t
}

func StringFromInterface(data interface{}) string {
	switch data.(type) {
	case string:
		return strings.TrimSpace(data.(string))
	case int:
		return fmt.Sprintf("%d", data.(int))
	case uint:
		return fmt.Sprintf("%d", data.(uint))
	case float64:
		return fmt.Sprintf("%.0f", data.(float64))
	default:
		return ""
	}
}

func BoolFromInterface(data ...interface{}) bool {
	switch data[0].(type) {
	case float64:
		return data[0].(float64) == 1
	case bool:
		return data[0].(bool)
	case nil:
		defaultVal := false
		dataLen := len(data)
		if dataLen > 1 {
			defaultVal = data[1].(bool)
		}
		return defaultVal
	default:
		return false
	}
}

func Urlencode(data map[string]string) string {
	params := url.Values{}
	for i, v := range data {
		params.Add(i, v)
	}
	return params.Encode()
}
