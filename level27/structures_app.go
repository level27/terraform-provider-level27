package level27

import (
	"encoding/json"
	"fmt"
	"strings"
)

// App represents a single App
type App struct {
	App struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Status       string `json:"status"`
		Organisation struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Reseller int    `json:"reseller"`
		} `json:"organisation"`
		DtExpires     string         `json:"dtExpires"`
		BillingStatus string         `json:"billingStatus"`
		Components    []AppComponent `json:"components"`
	} `json:"app"`
}

// AppRequest represents a single AppRequest
type AppRequest struct {
	Name         string
	Organisation string
}

func (ar AppRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", ar.Name)
	s += fmt.Sprintf("\"organisation\": %s", ar.Organisation)
	s += "}"
	return s
}

// AppComponent represents a single AppComponent
type AppComponent struct {
	Component struct {
		ID                     int    `json:"id"`
		Name                   string `json:"name"`
		Category               string `json:"category"`
		Appcomponenttype       string `json:"appcomponenttype"`
		BillableitemDetailID   int    `json:"billableitemDetailId"`
		Appcomponentparameters struct {
			User    string `json:"user"`
			Pass    string `json:"pass"`
			Path    string `json:"path"`
			Sshkeys []struct {
				Description    string `json:"description"`
				ID             int    `json:"id"`
				OrganisationID int    `json:"organisationId"`
				Owner          string `json:"owner"`
				Status         string `json:"status"`
				Type           string `json:"type"`
			} `json:"sshkeys"`
			Version     string `json:"version"`
			Size        string `json:"size"`
			Host        string `json:"host"`
			ExtraConfig string `json:"extraconfig"`
			Jailed      bool   `json:"jailed"`
		} `json:"appcomponentparameters"`
		AppcomponentparameterDescriptions struct {
			User        string `json:"user"`
			Pass        string `json:"pass"`
			Path        string `json:"path"`
			Sshkeys     string `json:"sshkeys"`
			Size        string `json:"size"`
			Host        string `json:"host"`
			ExtraConfig string `json:"extraconfig"`
			Jailed      string `json:"jailed"`
		} `json:"appcomponentparameterDescriptions"`
		Status string `json:"status"`
		App    struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		} `json:"app"`
		Organisation struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"organisation"`
		Systems []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"systems"`
		Systemgroup int `json:"systemgroup"`
		Provider    struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"provider"`
	} `json:"component"`
}

// AppComponentRequest represents a single AppRequest
type AppComponentRequest struct {
	Name        string
	Category    string
	Type        string
	System      string
	SystemGroup string
	Version     string
	Pass        string
	Path        string
	SSHKeys     []interface{}
	Size        string
	Host        string
	ExtraConfig string
	Jailed      bool
}

func (acr AppComponentRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", acr.Name)
	s += fmt.Sprintf("\"appcomponenttype\": \"%s\",", acr.Type)
	s += fmt.Sprintf("\"category\": \"%s\",", acr.Category)

	// mutually exclusive
	s += fmt.Sprintf("\"system\": \"%s\",", acr.System)
	s += fmt.Sprintf("\"systemgroup\": \"%s\",", acr.SystemGroup)

	// following depends on type
	switch acr.Type {
	case "asp":
		s += fmt.Sprintf("\"path\": \"%s\",", acr.Path)
		sshKeys, _ := json.Marshal(acr.SSHKeys)
		s += fmt.Sprintf("\"sshkeys\": %s,", sshKeys)
	case "php":
		s += fmt.Sprintf("\"version\": \"%s\",", acr.Version)
		s += fmt.Sprintf("\"extraconfig\": \"%s\",", acr.ExtraConfig)
		s += fmt.Sprintf("\"path\": \"%s\",", acr.Path)
		sshKeys, _ := json.Marshal(acr.SSHKeys)
		s += fmt.Sprintf("\"sshkeys\": %s,", sshKeys)
	case "sftp":
		s += fmt.Sprintf("\"jailed\": \"%v\",", acr.Jailed)
		sshKeys, _ := json.Marshal(acr.SSHKeys)
		s += fmt.Sprintf("\"sshkeys\": %s,", sshKeys)
	case "mssql":
		s += fmt.Sprintf("\"host\": \"%s\",", acr.Size)
	case "mysql":
		s += fmt.Sprintf("\"host\": \"%s\",", acr.Host)
	}

	s = strings.TrimSuffix(s, ",")
	s += "}"
	return s
}

// URL repr
type URL struct {
	URLs []struct {
		ID           int    `json:"id"`
		Content      string `json:"content"`
		HTTPS        bool   `json:"https"`
		Status       string `json:"status"`
		SslForce     bool   `json:"sslForce"`
		Appcomponent struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			Appcomponenttype string `json:"appcomponenttype"`
			Status           string `json:"status"`
		} `json:"appcomponent"`
		SslCertificate struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			SslStatus string `json:"sslStatus"`
			Status    string `json:"status"`
		} `json:"sslCertificate"`
		MatchingCertificates []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"matchingCertificates"`
		Type string `json:"type"`
	} `json:"urls"`
}

// URLRequest represents a single URLRequest
type URLRequest struct {
	Content        string
	SSLForce       bool
	SSLCertificate int
	HandleDNS      bool
}

func (ur URLRequest) String() string {
	s := "{"
	s += fmt.Sprintf("\"content\": \"%s\",", ur.Content)
	s += fmt.Sprintf("\"sslForce\": \"%t\",", ur.SSLForce)
	if ur.SSLCertificate != 0 {
		s += fmt.Sprintf("\"sslCertificate\": \"%d\",", ur.SSLCertificate)
	}
	s += fmt.Sprintf("\"handleDns\": \"%t\"", ur.HandleDNS)
	s += "}"
	return s
}
