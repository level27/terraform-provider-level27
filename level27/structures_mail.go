package level27

import (
	"fmt"
)

// Mailgroup represents a single Mailgroup
type Mailgroup struct {
	Mailgroup struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Status      string `json:"status"`
		Systemgroup struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"systemgroup"`
		Organisation struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Reseller interface{} `json:"reseller"`
		} `json:"organisation"`
		Domains       []interface{} `json:"domains"`
		DtExpires     string        `json:"dtExpires"`
		BillingStatus string        `json:"billingStatus"`
	} `json:"mailgroup"`
}

func (d Mailgroup) String() string {
	return "mailgroup"
}

// MailgroupRequest represents a single MailgroupRequest
type MailgroupRequest struct {
	Name                  string
	Action                string
	Domaintype            int
	Domaincontactlicensee string
	Organisation          string
	Handledns             bool
}

// LinkedDomain represents a single domain linked to a mailgroup
type LinkedDomain struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Domaintype struct {
		ID        int    `json:"id"`
		Extension string `json:"extension"`
	} `json:"domaintype"`
}

// Mailforwarder represents a single Mailforwarder
type Mailforwarder struct {
	Mailforwarder struct {
		ID          int      `json:"id"`
		Address     string   `json:"address"`
		Destination []string `json:"destination"`
		Status      string   `json:"status"`
		Mailgroup   struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"mailgroup"`
		Domain struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Domaintype struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"domaintype"`
		} `json:"domain"`
	} `json:"mailforwarder"`
}

// MailforwarderRequest represents a single MailforwarderRequest
type MailforwarderRequest struct {
	Address     string
	Destination string
}

// Mailbox represents a single Mailbox
type Mailbox struct {
	Mailbox struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Username   string `json:"username"`
		Status     string `json:"status"`
		OooEnabled bool   `json:"oooEnabled"`
		OooSubject string `json:"oooSubject"`
		OooText    string `json:"oooText"`
		Mailgroup  struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"mailgroup"`
		System struct {
			ID       int    `json:"id"`
			Fqdn     string `json:"fqdn"`
			Hostname string `json:"hostname"`
		} `json:"system"`
	} `json:"mailbox"`
}

// MailboxRequest represents a single MailboxRequest
type MailboxRequest struct {
	Name       string
	Password   string
	OooEnabled bool
	OooSubject string
	OooText    string
}

// MailboxAddress represents a single MailboxAddress
type MailboxAddress struct {
	MailboxAddress struct {
		ID      int    `json:"id"`
		Address string `json:"name"`
		Status  string `json:"status"`
	} `json:"mailboxAddress"`
}

// MailboxAddressRequest represents a single MailboxAddressRequest
type MailboxAddressRequest struct {
	Address string `json:"name"`
}

// MailboxAddresses represents an array of MailboxAddress
type MailboxAddresses struct {
	MailboxAddresses []struct {
		ID      int    `json:"id"`
		Address string `json:"address"`
		Status  string `json:"status"`
	} `json:"mailboxAddresses"`
}

func (d MailgroupRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", d.Name)
	s += fmt.Sprintf("\"organisation\": \"%s\",", d.Organisation)
	s += "\"type\": \"level27\""
	s += "}"
	return s
}

func (d LinkedDomain) String() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", d.Name)
	s += fmt.Sprintf("\"organisation\": \"%d\"", d.ID)
	s += "}"
	return s
}

func (d MailforwarderRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"address\": \"%s\",", d.Address)
	s += fmt.Sprintf("\"destination\": \"%s\"", d.Destination)
	s += "}"
	return s
}

func (d MailboxRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", d.Name)
	s += fmt.Sprintf("\"password\": \"%s\",", d.Password)
	s += fmt.Sprintf("\"oooEnabled\": %t,", d.OooEnabled)
	s += fmt.Sprintf("\"oooSubject\": \"%s\",", d.OooSubject)
	s += fmt.Sprintf("\"oooText\": \"%s\"", d.OooText)
	s += "}"
	return s
}

func (d MailboxAddressRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"address\": \"%s\"", d.Address)
	s += "}"
	return s
}
