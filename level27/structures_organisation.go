package level27

import (
	"fmt"
	"log"
)

// Organisation represents the organisation in CP4
type Organisation struct {
	Organisation struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		TaxNumber   string `json:"taxNumber"`
		MustPayTax  bool   `json:"mustPayTax"`
		Street      string `json:"street"`
		HouseNumber string `json:"houseNumber"`
		Zip         string `json:"zip"`
		City        string `json:"city"`
		Country     struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country,omitempty"`
		ResellerOrganisation struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			TaxNumber string `json:"taxNumber"`
		} `json:"resellerOrganisation,omitempty"`
		Users []struct {
			ID        int      `json:"id"`
			Username  string   `json:"username"`
			Email     string   `json:"email"`
			FirstName string   `json:"firstName"`
			LastName  string   `json:"lastName"`
			Roles     []string `json:"roles"`
		} `json:"users,omitempty"`
		RemarksToPrintInvoice string `json:"remarksToPrintInvoice"`
		UpdateEntitiesOnly    bool   `json:"updateEntitiesOnly"`
		PreventDeactivation   bool   `json:"preventDeactivation"`
	} `json:"organisation"`
}

// OrganisationRequest represents the request sent to CP4 API
type OrganisationRequest struct {
	Name                 string
	TaxNumber            string
	ResellerOrganisation string
}

func (o OrganisationRequest) String() string {
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
