package level27

import (
	"encoding/json"
	"fmt"
)

// Domain represents a single Domain
type Domain struct {
	Domain struct {
		ID                    int    `json:"id"`
		Name                  string `json:"name"`
		Fullname              string `json:"fullname"`
		TTL                   int    `json:"ttl"`
		EppCode               string `json:"eppCode"`
		Status                string `json:"status"`
		DnssecStatus          string `json:"dnssecStatus"`
		RegistrationIsHandled bool   `json:"registrationIsHandled"`
		Provider              string `json:"provider"`
		DNSIsHandled          bool   `json:"dnsIsHandled"`
		DtRegister            string `json:"dtRegister"`
		Nameserver1           string `json:"nameserver1"`
		Nameserver2           string `json:"nameserver2"`
		Nameserver3           string `json:"nameserver3"`
		Nameserver4           string `json:"nameserver4"`
		NameserverIP1         string `json:"nameserverIp1"`
		NameserverIP2         string `json:"nameserverIp2"`
		NameserverIP3         string `json:"nameserverIp3"`
		NameserverIP4         string `json:"nameserverIp4"`
		NameserverIpv61       string `json:"nameserverIpv61"`
		NameserverIpv62       string `json:"nameserverIpv62"`
		NameserverIpv63       string `json:"nameserverIpv63"`
		NameserverIpv64       string `json:"nameserverIpv64"`
		Organisation          struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Reseller int    `json:"reseller"`
		} `json:"organisation"`
		Domaintype struct {
			ID                                  int    `json:"id"`
			Name                                string `json:"name"`
			Extension                           string `json:"extension"`
			RenewPeriod                         int    `json:"renewPeriod"`
			TransferAutoLicensee                bool   `json:"transferAutoLicensee"`
			RequestIncomingTransferCodePossible bool   `json:"requestIncomingTransferCodePossible"`
			RequestOutgoingTransferCodePossible bool   `json:"requestOutgoingTransferCodePossible"`
			LicenseeChangePossible              bool   `json:"licenseeChangePossible"`
			DnssecSupported                     bool   `json:"dnssecSupported"`
		} `json:"domaintype"`
		DomaincontactLicensee struct {
			ID               int    `json:"id,omitempty"`
			FirstName        string `json:"firstName"`
			LastName         string `json:"lastName"`
			OrganisationName string `json:"organisationName"`
			Street           string `json:"street"`
			HouseNumber      string `json:"houseNumber"`
			Zip              string `json:"zip"`
			City             string `json:"city"`
			State            string `json:"state"`
			Phone            string `json:"phone"`
			Fax              string `json:"fax"`
			Email            string `json:"email"`
			TaxNumber        string `json:"taxNumber"`
			Status           int    `json:"status"`
			PassportNumber   string `json:"passportNumber"`
			SocialNumber     string `json:"socialNumber"`
			BirthStreet      string `json:"birthStreet"`
			BirthZip         string `json:"birthZip"`
			BirthCity        string `json:"birthCity"`
			BirthDate        string `json:"birthDate"`
			Gender           string `json:"gender"`
			Type             string `json:"type"`
			Country          struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"country"`
			Fullname string `json:"fullname"`
		} `json:"domaincontactLicensee"`
		DomaincontactOnsite interface{}   `json:"domaincontactOnsite"`
		Mailgroup           interface{}   `json:"mailgroup"`
		ExtraFields         []interface{} `json:"extraFields"`
		HandleMailDNS       interface{}   `json:"handleMailDns"`
		DtExpires           int           `json:"dtExpires"`
		BillingStatus       string        `json:"billingStatus"`
		ExternalInfo        interface{}   `json:"externalInfo"`
		Teams               []interface{} `json:"teams"`
		CountTeams          int           `json:"countTeams"`
	} `json:"domain"`
}

func (d Domain) String() string {
	return "domain"
}

// DomainProvider represents a single DomainProvider
type DomainProvider struct {
	Providers []struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		API             string `json:"api"`
		DNSSecSupported bool   `json:"dnsSecSupported"`
		Domaintypes     []struct {
			ID        int    `json:"id"`
			Extension string `json:"extension"`
		} `json:"domaintypes"`
	} `json:"providers"`
}

// DomainExtension represents a single DomainExtension
type DomainExtension struct {
	ID        int
	Extension string
}

// DomainRequest represents a single DomainRequest
type DomainRequest struct {
	Name                  string
	Action                string
	Domaintype            int
	Domaincontactlicensee string
	Organisation          string
	Handledns             bool
}

func (d DomainRequest) String() string {
	licensee := fmt.Sprintf("%s", d.Domaincontactlicensee)
	if d.Domaincontactlicensee == "" {
		licensee = "0"
	}

	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", d.Name)
	s += fmt.Sprintf("\"action\": \"%s\",", d.Action)
	s += fmt.Sprintf("\"domaintype\": %d,", d.Domaintype)
	s += fmt.Sprintf("\"domaincontactLicensee\": %s,", licensee)
	s += fmt.Sprintf("\"organisation\": \"%s\",", d.Organisation)
	s += fmt.Sprintf("\"handleDns\": %t", d.Handledns)
	s += "}"
	return s
}

// DomainRecord represents a single Domainrecord
type DomainRecord struct {
	Record struct {
		ID                 int    `json:"id"`
		Name               string `json:"name"`
		Content            string `json:"content"`
		Priority           int    `json:"priority"`
		Type               string `json:"type"`
		SystemHasNetworkIP struct {
			ID int `json:"id"`
		} `json:"systemHasNetworkIp"`
		URL            int `json:"url"`
		SslCertificate int `json:"sslCertificate"`
		Mailgroup      int `json:"mailgroup"`
	} `json:"record"`
}

// DomainRecordRequest represents a API reqest to Level27
type DomainRecordRequest struct {
	Name     string
	Type     string
	Content  string
	Priority int
}

func (d DomainRecordRequest) String() string {
	name := fmt.Sprintf("\"%s\"", d.Name)
	if d.Name == "" {
		name = "null"
	}

	s := "{"
	s += fmt.Sprintf("\"name\": %s,", name)
	s += fmt.Sprintf("\"type\": \"%s\",", d.Type)
	s += fmt.Sprintf("\"content\": \"%s\",", d.Content)
	s += fmt.Sprintf("\"priority\": %d", d.Priority)
	s += "}"
	return s
}

// DomainContact is an object to define domain contacts at Level27
type DomainContact struct {
	Domaincontact struct {
		ID               int    `json:"id"`
		FirstName        string `json:"firstName"`
		LastName         string `json:"lastName"`
		OrganisationName string `json:"organisationName"`
		Street           string `json:"street"`
		HouseNumber      string `json:"houseNumber"`
		Zip              string `json:"zip"`
		City             string `json:"city"`
		State            string `json:"state"`
		Phone            string `json:"phone"`
		Fax              string `json:"fax"`
		Email            string `json:"email"`
		TaxNumber        string `json:"taxNumber"`
		PassportNumber   string `json:"passportNumber"`
		SocialNumber     string `json:"socialNumber"`
		BirthStreet      string `json:"birthStreet"`
		BirthZip         string `json:"birthZip"`
		BirthCity        string `json:"birthCity"`
		BirthDate        string `json:"birthDate"`
		Gender           string `json:"gender"`
		Type             string `json:"type"`
		Country          struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
		Organisation struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"organisation"`
		Fullname string `json:"fullname"`
	} `json:"domaincontact"`
}

// DomainContactRequest is an object to define the request to create or modify a domain contact at Level27
type DomainContactRequest struct {
	Type             string `json:"type"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	OrganisationName string `json:"organisationName"`
	Street           string `json:"street"`
	HouseNumber      string `json:"houseNumber,omitempty"`
	Zip              string `json:"zip"`
	City             string `json:"city"`
	State            string `json:"state,omitempty"`
	Phone            string `json:"phone"`
	Fax              string `json:"fax,omitempty"`
	Email            string `json:"email"`
	TaxNumber        string `json:"taxNumber"`
	PassportNumber   string `json:"passportNumber,omitempty"`
	SocialNumber     string `json:"socialNumber,omitempty"`
	BirthStreet      string `json:"birthStreet,omitempty"`
	BirthZip         string `json:"birthZip,omitempty"`
	BirthCity        string `json:"birthCity,omitempty"`
	BirthDate        string `json:"birthDate,omitempty"`
	Gender           string `json:"gender,omitempty"`
	Country          string `json:"country"`
	Organisation     string `json:"organisation"`
}

func (d DomainContactRequest) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}
