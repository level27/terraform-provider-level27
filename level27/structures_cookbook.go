package level27

import (
	"fmt"
	"strings"
)

// Cookbook represents a single cookbook for a system
type Cookbook struct {
	Cookbook struct {
		ID                 int    `json:"id"`
		Cookbooktype       string `json:"cookbooktype"`
		Cookbookparameters map[string]struct {
			Value   interface{} `json:"value"`
			Default bool        `json:"default"`
		} `json:"cookbookparameters,omitempty"`
		CookbookparameterDescriptions map[string]string `json:"cookbookparameterDescriptions,omitempty"`
		PreviousCookbookparameters    string            `json:"previousCookbookparameters"`
		Status                        string            `json:"status"`
		System                        struct {
			ID   int    `json:"id"`
			Fqdn string `json:"fqdn"`
			Name string `json:"name"`
		} `json:"system"`
	} `json:"cookbook"`
}

func (c Cookbook) String() string {
	return "cookbook"
}

//CookbookRequest prepares an object to create a new cookbook
type CookbookRequest struct {
	Type               string
	Cookbookparameters map[string]interface{}
}

func (cr CookbookRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"cookbooktype\": \"%s\",", cr.Type)
	for key, elem := range cr.Cookbookparameters {
		if (cr.Type == "php" && key == "versions") || (cr.Type == "mysql" && key == "version") {
			s += fmt.Sprintf("\"%s\": [", key)
			for _, version := range elem.([]interface{}) {
				s += fmt.Sprintf("\"%v\",", version)
			}
			s = strings.TrimSuffix(s, ",")
			s += "],"
		} else {
			s += fmt.Sprintf("\"%s\": \"%v\",", key, elem)
		}
	}
	s = strings.TrimSuffix(s, ",")
	s += "}"
	return s
}
