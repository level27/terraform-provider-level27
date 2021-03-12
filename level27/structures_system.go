package level27

import (
	"fmt"
	"strconv"
	"strings"
)

// System represents a single System
type System struct {
	System struct {
		ID                int         `json:"id"`
		UID               string      `json:"uid"`
		Hostname          interface{} `json:"hostname"`
		Fqdn              string      `json:"fqdn"`
		Name              string      `json:"name"`
		Remarks           string      `json:"remarks"`
		Status            string      `json:"status"`
		RunningStatus     string      `json:"runningStatus"`
		CPU               interface{} `json:"cpu"`
		Memory            interface{} `json:"memory"`
		Disk              string      `json:"disk"`
		MonitoringEnabled bool        `json:"monitoringEnabled"`
		ManagementType    string      `json:"managementType"`
		Organisation      struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			Reseller interface{} `json:"reseller"`
		} `json:"organisation"`
		Systemimage struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			ExternalID string `json:"externalId"`
			OsID       int    `json:"osId"`
			OsName     string `json:"osName"`
			OsType     string `json:"osType"`
			OsVersion  string `json:"osVersion"`
		} `json:"systemimage"`
		OperatingsystemVersion struct {
			ID        int    `json:"id"`
			OsID      int    `json:"osId"`
			OsName    string `json:"osName"`
			OsType    string `json:"osType"`
			OsVersion string `json:"osVersion"`
		} `json:"operatingsystemVersion"`
		Provider struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			API  string `json:"api"`
		} `json:"provider"`
		SystemproviderConfiguration struct {
			ID          int    `json:"id"`
			ExternalID  string `json:"externalId"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"systemproviderConfiguration"`
		Region string `json:"region"`
		Zone   struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"zone"`
		Networks []struct {
			ID           int    `json:"id"`
			Mac          string `json:"mac"`
			NetworkID    int    `json:"networkId"`
			Name         string `json:"name"`
			NetIpv4      string `json:"netIpv4"`
			NetGatewayv4 string `json:"netGatewayv4"`
			NetMaskv4    int    `json:"netMaskv4"`
			NetIpv6      string `json:"netIpv6"`
			NetGatewayv6 string `json:"netGatewayv6"`
			NetMaskv6    int    `json:"netMaskv6"`
			NetPublic    bool   `json:"netPublic"`
			NetCustomer  bool   `json:"netCustomer"`
			NetInternal  bool   `json:"netInternal"`
			Ips          []struct {
				ID         int    `json:"id"`
				PublicIpv4 string `json:"publicIpv4"`
				Ipv4       string `json:"ipv4"`
				PublicIpv6 string `json:"publicIpv6"`
				Ipv6       string `json:"ipv6"`
				Hostname   string `json:"hostname"`
			} `json:"ips"`
			Destinationv4 []string `json:"destinationv4"`
			Destinationv6 []string `json:"destinationv6"`
		} `json:"networks"`
		PublicNetworking bool `json:"publicNetworking"`
		StatsSummary     struct {
			Diskspace struct {
				Unit  string      `json:"unit"`
				Value interface{} `json:"value"`
				Max   interface{} `json:"max"`
			} `json:"diskspace"`
			Memory struct {
				Unit  string      `json:"unit"`
				Value interface{} `json:"value"`
				Max   interface{} `json:"max"`
			} `json:"memory"`
			CPU struct {
				Unit  string      `json:"unit"`
				Value interface{} `json:"value"`
				Max   int         `json:"max"`
			} `json:"cpu"`
		} `json:"statsSummary"`
		Group                  interface{}   `json:"group"`
		BillingStatus          string        `json:"billingStatus"`
		DtExpires              int           `json:"dtExpires"`
		Cookbooks              []interface{} `json:"cookbooks"`
		Volumes                []interface{} `json:"volumes"`
		InstallSecurityUpdates int           `json:"installSecurityUpdates"`
	} `json:"system"`
}

//SystemRequest prepares an object to create a new system
type SystemRequest struct {
	CustomerFqdn                string
	Name                        string
	Remarks                     string
	CPU                         int
	Memory                      int
	Disk                        int
	Systemimage                 string
	Organisation                string
	SystemproviderConfiguration string
	Zone                        string
	InstallSecurityUpdates      int
	AutoTeams                   string
	ExternalInfo                string
	AutoNetworks                []interface{}
}

// SystemProviders provides a list of all system providers with their images
type SystemProviders struct {
	Providers []struct {
		ID                 int           `json:"id"`
		Name               string        `json:"name"`
		API                string        `json:"api"`
		AdvancedNetworking bool          `json:"advancedNetworking"`
		Icon               string        `json:"icon"`
		Images             []SystemImage `json:"images"`
	} `json:"providers"`
}

// SystemImage provides a single instance of a system image
type SystemImage struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"externalId"`
	Region     struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"region"`
	OperatingSystemID      int    `json:"operatingSystemID"`
	OperatingSystem        string `json:"operatingSystem"`
	OperatingSystemVersion struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"operatingSystemVersion"`
}

// Regions provides a list of all regions
type Regions struct {
	Regions []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
		Zones []Zone `json:"zones"`
	} `json:"regions"`
}

// Zone represents a zone for a System to live in
type Zone struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

// ProviderConfigurations provides a list of provider configurations
type ProviderConfigurations struct {
	ProviderConfigurations []struct {
		ID             int    `json:"id"`
		MinCPU         int    `json:"minCpu"`
		MaxCPU         int    `json:"maxCpu"`
		MinMemory      string `json:"minMemory"`
		MaxMemory      string `json:"maxMemory"`
		MinDisk        int    `json:"minDisk"`
		MaxDisk        int    `json:"maxDisk"`
		Status         int    `json:"status"`
		ExternalID     string `json:"externalId"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Systemprovider struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"systemprovider"`
	} `json:"providerConfigurations"`
}

func (s System) String() string {
	output := fmt.Sprintf("\"name\": \"%s\"\n", s.System.Name)
	output += fmt.Sprintf("\"status\": \"%s\"\n", s.System.Status)
	output += fmt.Sprintf("\"running\": \"%s\"\n", s.System.RunningStatus)
	return output
}

func (p SystemProviders) getSystemProviders() []string {
	names := []string{}
	for _, provider := range p.Providers {
		names = append(names, strconv.Itoa(provider.ID)+"-"+provider.Name)
	}
	return names
}

func (p SystemProviders) getSystemImagesByProviderID(providerID int) []SystemImage {
	images := []SystemImage{}
	for _, provider := range p.Providers {
		if providerID == provider.ID {
			for _, image := range provider.Images {
				images = append(images, image)
			}
		}
	}
	return images
}

func (p SystemProviders) getSystemImagesByProviderName(providerName string) []SystemImage {
	images := []SystemImage{}
	for _, provider := range p.Providers {
		if providerName == provider.Name {
			for _, image := range provider.Images {
				images = append(images, image)
			}
		}
	}
	return images
}

func (sr SystemRequest) String() string {
	networks := ""
	for _, network := range sr.AutoNetworks {
		networks += fmt.Sprintf("\"%v\",", network)
	}
	networks = strings.TrimSuffix(networks, ",")

	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", sr.Name)
	s += fmt.Sprintf("\"customerFqdn\": \"%s\",", sr.CustomerFqdn)
	s += fmt.Sprintf("\"cpu\": \"%d\",", sr.CPU)
	s += fmt.Sprintf("\"memory\": \"%d\",", sr.Memory)
	s += fmt.Sprintf("\"disk\": \"%d\",", sr.Disk)
	s += fmt.Sprintf("\"systemimage\": %s,", sr.Systemimage)
	s += fmt.Sprintf("\"organisation\": %s,", sr.Organisation)
	s += fmt.Sprintf("\"systemproviderConfiguration\": %s,", sr.SystemproviderConfiguration)
	s += fmt.Sprintf("\"zone\": %s,", sr.Zone)
	s += fmt.Sprintf("\"autoNetworks\": \"%s\"", networks)

	return s
}

// UpdateString creates data to send in PUT request for System
func (sr SystemRequest) UpdateString() string {
	s := "{"
	s += fmt.Sprintf("\"name\": \"%s\",", sr.Name)
	s += fmt.Sprintf("\"customerFqdn\": \"%s\",", sr.CustomerFqdn)
	s += fmt.Sprintf("\"disk\": \"%d\",", sr.Disk)
	s += fmt.Sprintf("\"systemimage\": %s,", sr.Systemimage)
	s += fmt.Sprintf("\"organisation\": %s,", sr.Organisation)
	s += fmt.Sprintf("\"systemproviderConfiguration\": %s,", sr.SystemproviderConfiguration)
	s += fmt.Sprintf("\"zone\": %s", sr.Zone)
	s += "}"
	return s
}
