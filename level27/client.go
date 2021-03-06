package level27

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Client defines the API Client structure
type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

// NewAPIClient creates a client for doing the API calls
func NewAPIClient(uri string, apiKey string) *Client {
	return &Client{
		BaseURL: uri,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  struct {
		Children struct {
			Content struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"content,omitempty"`
			SSLForce struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslForce,omitempty"`
			SSLCertificate struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"sslCertificate,omitempty"`
			HandleDNS struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"handleDns,omitempty"`
			Authentication struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"authentication,omitempty"`
			Appcomponent struct {
				Errors []string `json:"errors,omitempty"`
			} `json:"appcomponent,omitempty"`
		} `json:"children"`
	} `json:"errors"`
}

func (er errorResponse) String() string {
	s := fmt.Sprintf("%s\n", er.Message)
	fields := reflect.TypeOf(er.Errors.Children)
	values := reflect.ValueOf(er.Errors.Children)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		s += fmt.Sprintf("%v = %v\n", field.Name, value)
	}
	return s
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (c *Client) sendRequest(method string, endpoint string, data interface{}) ([]byte, error) {

	log.Printf("client.go: Send %s request > %s/%s", method, c.BaseURL, endpoint)
	log.Printf("client.go: Body: %v", data)

	reqData := bytes.NewBuffer([]byte(nil))
	if data != nil {
		reqData = bytes.NewBuffer([]byte(fmt.Sprintf("%v", data)))
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.BaseURL, endpoint), reqData)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if method == "UPDATE" && res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(res.Body)
		jsonParsed, err := gabs.ParseJSON(body)
		if err != nil {
			return nil, err
		}

		log.Printf("client.go: ERROR: %v", jsonParsed)
		for key, child := range jsonParsed.Search("errors", "children").ChildrenMap() {
			if child.Data().(map[string]interface{})["errors"] != nil {
				errorMessages := child.Data().(map[string]interface{})["errors"].([]interface{})
				if len(errorMessages) > 0 {
					for _, err := range errorMessages {
						log.Printf("Key=>%v, Value=>%v\n", key, err)
						return nil, fmt.Errorf("%v : %v", key, err)
					}
				}
			}
		}

		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, errors.New(errRes.String())
		}

		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	return body, nil
}

func (c *Client) invokeAPI(method string, endpoint string, data interface{}, result interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	body, err := c.sendRequest(method, endpoint, data)

	if err != nil {
		log.Printf("client.go: API error - %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Level27 API error",
			Detail:   fmt.Sprintf("Error returned from Level27 API\n%v", err.Error()),
		})
	}

	if method == "PUT" || method == "DELETE" || method == "PATCH" {
		return diags
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Level27 API error",
			Detail:   fmt.Sprintf("Unable to unmarshall the body"),
		})
	}

	return diags
}
