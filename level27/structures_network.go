package level27

// Network represents a single network
type Network struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Remarks         string      `json:"remarks"`
	Status          string      `json:"status"`
	Vlan            int         `json:"vlan"`
	Ipv4            string      `json:"ipv4"`
	Netmaskv4       int         `json:"netmaskv4"`
	Gatewayv4       string      `json:"gatewayv4"`
	Ipv6            string      `json:"ipv6"`
	Netmaskv6       int         `json:"netmaskv6"`
	Gatewayv6       string      `json:"gatewayv6"`
	Public          bool        `json:"public"`
	Customer        bool        `json:"customer"`
	Internal        bool        `json:"internal"`
	PublicIP4Native interface{} `json:"publicIp4Native"`
	PublicIP6Native interface{} `json:"publicIp6Native"`
	Full            interface{} `json:"full"`
	Systemgroup     interface{} `json:"systemgroup"`
	Organisation    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"organisation"`
	Zone struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Region struct {
			ID             int `json:"id"`
			Systemprovider struct {
				ID                 int    `json:"id"`
				API                string `json:"api"`
				Name               string `json:"name"`
				AdvancedNetworking bool   `json:"advancedNetworking"`
			} `json:"systemprovider"`
		} `json:"region"`
	} `json:"zone"`
}

// Networks represents an array of network
type Networks struct {
	Networks []Network `json:"networks"`
}
