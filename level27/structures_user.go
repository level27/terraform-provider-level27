package level27

import (
	"fmt"
	"log"
)

// User represents a single user in CP4
type User struct {
	User struct {
		ID           int `json:"id"`
		Invitation   bool
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		FirstName    string   `json:"firstName"`
		LastName     string   `json:"lastName"`
		Roles        []string `json:"roles"`
		Status       string   `json:"status"`
		Language     string   `json:"language"`
		Organisation struct {
			ID          int         `json:"id"`
			Name        string      `json:"name"`
			Street      string      `json:"street"`
			HouseNumber string      `json:"houseNumber"`
			Zip         string      `json:"zip"`
			City        string      `json:"city"`
			Reseller    interface{} `json:"reseller"`
		} `json:"organisation"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
		Fullname string `json:"fullname"`
		Teams    []struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			AutoAddSshkey bool   `json:"autoAddSshkey"`
		} `json:"teams"`
	} `json:"user"`
}

// UserRequest represents the request sent to CP4 API
type UserRequest struct {
	Name                 string
	TaxNumber            string
	ResellerOrganisation string
}

func (o UserRequest) String() string {
	log.Printf("start request")
	reseller := fmt.Sprintf("\"%s\"", o.ResellerOrganisation)
	if len(o.ResellerOrganisation) == 0 {
		reseller = "null"
	}

	log.Printf("%s", reseller)

	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", o.Name)
	s += fmt.Sprintf("\"taxNumber\": \"%s\",", o.TaxNumber)
	s += fmt.Sprintf("\"resellerOrganisation\": %s", reseller)
	s += "}"
	return s
}
