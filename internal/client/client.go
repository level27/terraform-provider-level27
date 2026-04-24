// Copyright (c) Level27 NV
// SPDX-License-Identifier: MPL-2.0

// Package client provides a Go HTTP client for the Level27 CP4 API.
// API base URL: https://api.level27.eu/v1/
// Authentication: Authorization: <api-key> header.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.level27.eu/v1"

// jsonInt is an integer that the Level27 API sometimes returns as a quoted
// decimal string (e.g. "60.0" for disk size).  It unmarshals both forms.
type jsonInt int

func (j *jsonInt) UnmarshalJSON(data []byte) error {
	// plain JSON number
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*j = jsonInt(n)
		return nil
	}
	// quoted string like "60.0"
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("jsonInt: expected number or string, got %s", data)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("jsonInt: cannot parse %q as number: %w", s, err)
	}
	*j = jsonInt(int(f))
	return nil
}

// Client is the Level27 API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// New creates a new Level27 API client.
func New(apiKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// ---------------------------------------------------------------------------
// Low-level HTTP helpers
// ---------------------------------------------------------------------------

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	u := c.baseURL + path

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("level27 client: marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("level27 client: http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("level27 client: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, raw)
	}

	if v != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, v); err != nil {
			return fmt.Errorf("level27 client: unmarshal response: %w", err)
		}
	}
	return nil
}

// APIError represents an error returned by the Level27 API.
type APIError struct {
	StatusCode int
	Message    string
	Details    map[string][]string
	RawBody    string
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("level27 api error %d: %s", e.StatusCode, e.Message)
	if len(e.Details) > 0 {
		for field, errs := range e.Details {
			msg += fmt.Sprintf("\n  - %s: %s", field, strings.Join(errs, ", "))
		}
	}
	if e.RawBody != "" {
		msg += fmt.Sprintf("\n  raw: %s", e.RawBody)
	}
	return msg
}

func parseAPIError(statusCode int, body []byte) *APIError {
	var payload struct {
		Message string              `json:"message"`
		Errors  map[string][]string `json:"errors"`
	}
	_ = json.Unmarshal(body, &payload)
	if payload.Message == "" {
		payload.Message = string(body)
	}
	return &APIError{StatusCode: statusCode, Message: payload.Message, Details: payload.Errors, RawBody: string(body)}
}

// IsNotFound returns true when the API returned HTTP 404.
func IsNotFound(err error) bool {
	if e, ok := err.(*APIError); ok {
		return e.StatusCode == 404
	}
	return false
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

// Ref is a minimal object reference (id + name), used throughout the API.
type Ref struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PaginationParams are common query parameters for list endpoints.
type PaginationParams struct {
	Limit  int
	Offset int
	Filter string
}

func (p PaginationParams) apply(q url.Values) {
	if p.Limit > 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	} else {
		q.Set("limit", "100") // override the tiny API default of 5
	}
	if p.Offset > 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Filter != "" {
		q.Set("filter", p.Filter)
	}
}

// ---------------------------------------------------------------------------
// System
// ---------------------------------------------------------------------------

// SystemService represents a single service installed on a system.
type SystemService struct {
	ID                int                         `json:"id"`
	ServiceType       string                      `json:"cookbooktype"`
	ServiceParameters map[string]ServiceParameter `json:"cookbookparameters"`
}

// ServiceParameter holds the value for a single service parameter.
// The value can be a string, number, bool, or array of strings.
type ServiceParameter struct {
	Value json.RawMessage `json:"value"`
}

// Versions returns the list of versions for this service parameter.
// The API returns either a string or []string.
func (cp ServiceParameter) Versions() []string {
	var arr []string
	if err := json.Unmarshal(cp.Value, &arr); err == nil {
		return arr
	}
	var s string
	if err := json.Unmarshal(cp.Value, &s); err == nil && s != "" {
		return []string{s}
	}
	return nil
}

// System represents a Level27 system (server).
type System struct {
	ID                          int             `json:"id"`
	Name                        string          `json:"name"`
	CustomerFqdn                string          `json:"customerFqdn"`
	Type                        string          `json:"type"`
	Status                      string          `json:"status"`
	StatusCategory              string          `json:"statusCategory"`
	CPU                         int             `json:"cpu"`
	Disk                        jsonInt         `json:"disk"`
	Memory                      int             `json:"memory"`
	ManagementType              string          `json:"managementType"`
	ExternalInfo                string          `json:"externalInfo"`
	Organisation                *Ref            `json:"organisation"`
	Zone                        *Ref            `json:"zone"`
	Systemimage                 *Ref            `json:"systemimage"`
	SystemproviderConfiguration *Ref            `json:"systemproviderConfiguration"`
	Parentsystem                *Ref            `json:"parentsystem"`
	Services                    []SystemService `json:"cookbooks"`
}

// CreateSystemRequest is the body for POST /systems.
type CreateSystemRequest struct {
	Name                        string  `json:"name"`
	CustomerFqdn                string  `json:"customerFqdn,omitempty"`
	Type                        string  `json:"type"`
	Organisation                int     `json:"organisation"`
	Systemimage                 int     `json:"systemimage"`
	SystemproviderConfiguration int     `json:"systemproviderConfiguration"`
	Zone                        int     `json:"zone"`
	CPU                         int     `json:"cpu"`
	Disk                        int     `json:"disk"`
	Memory                      int     `json:"memory"`
	ManagementType              string  `json:"managementType,omitempty"`
	Period                      int     `json:"period,omitempty"`
	ExternalInfo                *string `json:"externalInfo"`
	Parentsystem                *int    `json:"parentsystem"`
}

// UpdateSystemRequest is the body for PUT /systems/{id}.
type UpdateSystemRequest struct {
	Name                        string  `json:"name"`
	CustomerFqdn                string  `json:"customerFqdn,omitempty"`
	Type                        string  `json:"type"`
	Organisation                int     `json:"organisation"`
	Systemimage                 int     `json:"systemimage"`
	SystemproviderConfiguration int     `json:"systemproviderConfiguration"`
	Zone                        int     `json:"zone"`
	CPU                         int     `json:"cpu"`
	Disk                        int     `json:"disk"`
	Memory                      int     `json:"memory"`
	ManagementType              string  `json:"managementType,omitempty"`
	ExternalInfo                *string `json:"externalInfo"`
}

// GetSystem calls GET /systems/{id}.
func (c *Client) GetSystem(ctx context.Context, id int) (*System, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/systems/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		System System `json:"system"`
	}
	return &out.System, c.do(req, &out)
}

// CreateSystem calls POST /systems.
func (c *Client) CreateSystem(ctx context.Context, body CreateSystemRequest) (*System, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/systems", body)
	if err != nil {
		return nil, err
	}
	var out struct {
		System System `json:"system"`
	}
	return &out.System, c.do(req, &out)
}

// UpdateSystem calls PUT /systems/{id}.
func (c *Client) UpdateSystem(ctx context.Context, id int, body UpdateSystemRequest) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/systems/%d", id), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// DeleteSystem calls DELETE /systems/{id}.
func (c *Client) DeleteSystem(ctx context.Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/systems/%d", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// PerformSystemAction calls POST /systems/{id}/actions with the given action name.
// Common actions: "autoInstall", "start", "stop", "reboot".
func (c *Client) PerformSystemAction(ctx context.Context, systemID int, action string) error {
	body := struct {
		Action string `json:"action"`
	}{Action: action}
	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/systems/%d/actions", systemID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// WaitForSystemStatus polls GET /systems/{id} until the status is no longer
// transitional. Returns the final system or an error.
func (c *Client) WaitForSystemStatus(ctx context.Context, id int) (*System, error) {
	transitional := map[string]bool{
		"to_create":  true,
		"creating":   true,
		"to_update":  true,
		"updating":   true,
		"to_delete":  true,
		"deleting":   true,
		"to_install": true,
		"installing": true,
	}
	consecutiveErrors := 0
	const maxConsecutiveErrors = 5
	for {
		sys, err := c.GetSystem(ctx, id)
		if err != nil {
			// Tolerate transient network errors (timeouts, connection resets) for
			// up to maxConsecutiveErrors attempts before giving up.
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			consecutiveErrors++
			if consecutiveErrors > maxConsecutiveErrors {
				return nil, fmt.Errorf("too many consecutive errors waiting for system status: %w", err)
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(15 * time.Second):
			}
			continue
		}
		consecutiveErrors = 0
		if !transitional[sys.Status] {
			if sys.StatusCategory == "red" {
				return nil, fmt.Errorf("system reached error status: %s", sys.Status)
			}
			return sys, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

// FindServiceVersion looks up the highest available version for the given
// service type on this system. Returns "" if not found.
// The API stores versions under either "version" or "versions" key.
func (s *System) FindServiceVersion(serviceType string) string {
	for _, svc := range s.Services {
		if svc.ServiceType != serviceType {
			continue
		}
		var all []string
		for _, key := range []string{"version", "versions"} {
			if vp, ok := svc.ServiceParameters[key]; ok {
				all = append(all, vp.Versions()...)
			}
		}
		if len(all) == 0 {
			return ""
		}
		sort.Slice(all, func(i, j int) bool {
			return compareVersions(all[i], all[j]) > 0
		})
		return all[0]
	}
	return ""
}

// compareVersions compares two version strings (e.g. "8.2", "7.4").
// Returns positive if a > b, negative if a < b, 0 if equal.
func compareVersions(a, b string) int {
	partsA := strings.SplitN(a, ".", 3)
	partsB := strings.SplitN(b, ".", 3)
	for len(partsA) < 3 {
		partsA = append(partsA, "0")
	}
	for len(partsB) < 3 {
		partsB = append(partsB, "0")
	}
	for i := 0; i < 3; i++ {
		na, _ := strconv.Atoi(partsA[i])
		nb, _ := strconv.Atoi(partsB[i])
		if na != nb {
			return na - nb
		}
	}
	return 0
}

// ListSystems calls GET /systems.
func (c *Client) ListSystems(ctx context.Context, p PaginationParams) ([]System, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/systems?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Systems []System `json:"systems"`
	}
	return out.Systems, c.do(req, &out)
}

// ---------------------------------------------------------------------------
// Organisation
// ---------------------------------------------------------------------------

// Organisation represents a Level27 organisation.
type Organisation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ListOrganisations calls GET /organisations.
func (c *Client) ListOrganisations(ctx context.Context, p PaginationParams) ([]Organisation, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/organisations?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Organisations []Organisation `json:"organisations"`
	}
	return out.Organisations, c.do(req, &out)
}

// ---------------------------------------------------------------------------
// SystemProvider / SystemImage / SystemProviderConfiguration / SystemZone
// ---------------------------------------------------------------------------

// SystemProvider represents a Level27 system provider (hypervisor backend).
// Images are embedded directly in the API response.
type SystemProvider struct {
	ID     int           `json:"id"`
	Name   string        `json:"name"`
	Images []SystemImage `json:"images"`
}

// SystemImage represents a bootable OS image available on a provider.
type SystemImage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SystemProviderConfiguration represents a hardware profile.
type SystemProviderConfiguration struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ExternalID     string `json:"externalId"`
	MinCPU         int    `json:"minCpu"`
	MaxCPU         int    `json:"maxCpu"`
	MinMemory      string `json:"minMemory"`
	MaxMemory      string `json:"maxMemory"`
	MinDisk        int    `json:"minDisk"`
	MaxDisk        int    `json:"maxDisk"`
	Systemprovider *Ref   `json:"systemprovider"`
}

// SystemRegion represents a datacenter region containing zones.
type SystemRegion struct {
	ID    int          `json:"id"`
	Name  string       `json:"name"`
	Zones []SystemZone `json:"zones"`
}

// SystemZone represents a datacenter zone inside a region.
type SystemZone struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	RegionName string `json:"-"` // populated client-side from parent region
}

// ListSystemProviders calls GET /systems/providers.
func (c *Client) ListSystemProviders(ctx context.Context, p PaginationParams) ([]SystemProvider, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/systems/providers?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Providers []SystemProvider `json:"providers"`
	}
	return out.Providers, c.do(req, &out)
}

// ListSystemImages returns all images for the given provider.
// Images are embedded in the GET /systems/providers response.
func (c *Client) ListSystemImages(ctx context.Context, providerID int, p PaginationParams) ([]SystemImage, error) {
	providers, err := c.ListSystemProviders(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, pr := range providers {
		if pr.ID == providerID {
			return pr.Images, nil
		}
	}
	return nil, fmt.Errorf("provider %d not found", providerID)
}

// ListSystemProviderConfigurations calls GET /systems/provider/configurations.
// Returns all configurations across all providers.
func (c *Client) ListSystemProviderConfigurations(ctx context.Context, p PaginationParams) ([]SystemProviderConfiguration, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/systems/provider/configurations?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Configurations []SystemProviderConfiguration `json:"providerConfigurations"`
	}
	return out.Configurations, c.do(req, &out)
}

// ListSystemRegions calls GET /systems/regions and returns regions with nested zones.
func (c *Client) ListSystemRegions(ctx context.Context, p PaginationParams) ([]SystemRegion, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/systems/regions?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Regions []SystemRegion `json:"regions"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	// Backfill RegionName for convenience.
	for i, r := range out.Regions {
		for j := range out.Regions[i].Zones {
			out.Regions[i].Zones[j].RegionName = r.Name
		}
	}
	return out.Regions, nil
}

// ManagementType represents a system management level available for an organisation.
type ManagementType struct {
	// Value is the terraform management_type value (e.g. "infra_plus").
	Value       string
	Description string
	// DefaultPrice is the monthly price in euro cents (period=1, default=true).
	DefaultPrice string
}

// ListManagementTypes calls GET /systems/priceproposal/organisation/{orgID}
// and returns the available management types for that organisation.
func (c *Client) ListManagementTypes(ctx context.Context, orgID int) ([]ManagementType, error) {
	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/systems/priceproposal/organisation/%d", orgID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Products map[string][]struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Prices      []struct {
				Period  int    `json:"period"`
				Price   string `json:"price"`
				Default bool   `json:"default"`
			} `json:"prices"`
		} `json:"products"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	const prefix = "system_management_"
	var result []ManagementType
	seen := map[string]bool{}
	for _, products := range out.Products {
		for _, p := range products {
			if !strings.HasPrefix(p.ID, prefix) {
				continue
			}
			value := strings.TrimPrefix(p.ID, prefix)
			if seen[value] {
				continue
			}
			seen[value] = true
			defaultPrice := ""
			for _, pr := range p.Prices {
				if pr.Period == 1 && pr.Default {
					defaultPrice = pr.Price
					break
				}
			}
			result = append(result, ManagementType{
				Value:        value,
				Description:  p.Description,
				DefaultPrice: defaultPrice,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Value < result[j].Value })
	return result, nil
}

// Network represents a Level27 network.
type Network struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	IPv4        string `json:"ipv4"`
	NetmaskV4   int    `json:"netmaskv4"`
	IPv6        string `json:"ipv6"`
	NetmaskV6   int    `json:"netmaskv6"`
	Public      bool   `json:"public"`
	Customer    bool   `json:"customer"`
	Internal    bool   `json:"internal"`
}

// ListNetworks calls GET /networks?type=<networkType>.
// networkType is one of: "public", "customer", "internal".
func (c *Client) ListNetworks(ctx context.Context, networkType string, p PaginationParams) ([]Network, error) {
	q := url.Values{}
	p.apply(q)
	if networkType != "" {
		q.Set("type", networkType)
	}
	req, err := c.newRequest(ctx, http.MethodGet, "/networks?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Networks []Network `json:"networks"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return out.Networks, nil
}

// GetNetworkByName searches all network types for a network with the given name.
func (c *Client) GetNetworkByName(ctx context.Context, name string) (*Network, error) {
	p := PaginationParams{Limit: 10000}
	for _, t := range []string{"public", "customer", "internal"} {
		nets, err := c.ListNetworks(ctx, t, p)
		if err != nil {
			return nil, fmt.Errorf("searching %s networks: %w", t, err)
		}
		for i, n := range nets {
			if n.Name == name {
				return &nets[i], nil
			}
		}
	}
	return nil, fmt.Errorf("network with name %q not found", name)
}

// SystemHasNetwork represents the link between a system and a network.
type SystemHasNetwork struct {
	ID      int     `json:"id"`
	MAC     string  `json:"mac"`
	Status  string  `json:"status"`
	Network Network `json:"network"`
}

// ListSystemNetworks calls GET /systems/{id}/networks.
func (c *Client) ListSystemNetworks(ctx context.Context, systemID int) ([]SystemHasNetwork, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/systems/%d/networks", systemID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		SystemHasNetworks []SystemHasNetwork `json:"systemHasNetworks"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return out.SystemHasNetworks, nil
}

// AddSystemNetwork calls POST /systems/{id}/networks with body {"network": networkID}.
func (c *Client) AddSystemNetwork(ctx context.Context, systemID, networkID int) (*SystemHasNetwork, error) {
	body := struct {
		Network int `json:"network"`
	}{Network: networkID}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/systems/%d/networks", systemID), body)
	if err != nil {
		return nil, err
	}
	var out struct {
		SystemHasNetwork SystemHasNetwork `json:"systemHasNetwork"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.SystemHasNetwork, nil
}

// RemoveSystemNetwork calls DELETE /systems/{systemID}/networks/{systemHasNetworkID}.
func (c *Client) RemoveSystemNetwork(ctx context.Context, systemID, systemHasNetworkID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/systems/%d/networks/%d", systemID, systemHasNetworkID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// SystemHasNetworkIP represents an IP address assigned to a network interface.
type SystemHasNetworkIP struct {
	ID       int    `json:"id"`
	IPv4     string `json:"ipv4"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

// LocateNetworkIP calls GET /networks/{id}/locate and returns the first suggested free IPv4.
func (c *Client) LocateNetworkIP(ctx context.Context, networkID int) (string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/networks/%d/locate", networkID), nil)
	if err != nil {
		return "", err
	}
	var out struct {
		IPv4 []string `json:"ipv4"`
	}
	if err := c.do(req, &out); err != nil {
		return "", err
	}
	if len(out.IPv4) == 0 {
		return "", fmt.Errorf("no free IPv4 addresses available in network %d", networkID)
	}
	return out.IPv4[0], nil
}

// ListSystemNetworkIPs calls GET /systems/{id}/networks/{shnID}/ips.
func (c *Client) ListSystemNetworkIPs(ctx context.Context, systemID, shnID int) ([]SystemHasNetworkIP, error) {
	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/systems/%d/networks/%d/ips", systemID, shnID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		SystemHasNetworkIps []SystemHasNetworkIP `json:"systemHasNetworkIps"`
	}
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return out.SystemHasNetworkIps, nil
}

// AddSystemNetworkIP calls POST /systems/{id}/networks/{shnID}/ips with the given IPv4.
func (c *Client) AddSystemNetworkIP(ctx context.Context, systemID, shnID int, ipv4 string) error {
	body := struct {
		IPv4       string  `json:"ipv4"`
		PublicIPv4 *string `json:"publicIpv4"`
		IPv6       *string `json:"ipv6"`
		PublicIPv6 *string `json:"publicIpv6"`
		Hostname   string  `json:"hostname"`
	}{IPv4: ipv4}
	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/systems/%d/networks/%d/ips", systemID, shnID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// RemoveSystemNetworkIP calls DELETE /systems/{id}/networks/{shnID}/ips/{ipID}.
func (c *Client) RemoveSystemNetworkIP(ctx context.Context, systemID, shnID, ipID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/systems/%d/networks/%d/ips/%d", systemID, shnID, ipID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---------------------------------------------------------------------------
// App
// ---------------------------------------------------------------------------

// App represents a Level27 app/project.
type App struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	StatusCategory    string `json:"statusCategory"`
	Type              string `json:"type"`
	HostingType       string `json:"hostingType"`
	DtCreated         string `json:"dtCreated"`
	DtUpdated         string `json:"dtUpdated"`
	CustomPackageName string `json:"customPackageName"`
	BillingStatus     string `json:"billingStatus"`
	Organisation      *Ref   `json:"organisation"`
	Systemgroup       *Ref   `json:"systemgroup"`
}

// CreateAppRequest is the body for POST /apps.
type CreateAppRequest struct {
	Name              string  `json:"name"`
	Organisation      int     `json:"organisation"`
	CustomPackage     *int    `json:"customPackage,omitempty"`
	CustomPackageName *string `json:"customPackageName,omitempty"`
	AutoTeams         string  `json:"autoTeams,omitempty"`
	AutoUpgrades      string  `json:"autoUpgrades,omitempty"`
	ExternalInfo      string  `json:"externalInfo,omitempty"`
}

// UpdateAppRequest is the body for PUT /apps/{appId}.
type UpdateAppRequest struct {
	Name              string  `json:"name"`
	Organisation      int     `json:"organisation"`
	CustomPackage     *int    `json:"customPackage,omitempty"`
	CustomPackageName *string `json:"customPackageName,omitempty"`
	AutoTeams         string  `json:"autoTeams,omitempty"`
	AutoUpgrades      string  `json:"autoUpgrades,omitempty"`
	ExternalInfo      string  `json:"externalInfo,omitempty"`
}

// ListApps calls GET /apps.
func (c *Client) ListApps(ctx context.Context, p PaginationParams) ([]App, error) {
	q := url.Values{}
	p.apply(q)
	req, err := c.newRequest(ctx, http.MethodGet, "/apps?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Apps []App `json:"apps"`
	}
	return out.Apps, c.do(req, &out)
}

// CreateApp calls POST /apps.
func (c *Client) CreateApp(ctx context.Context, body CreateAppRequest) (*App, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/apps", body)
	if err != nil {
		return nil, err
	}
	var out struct {
		App App `json:"app"`
	}
	return &out.App, c.do(req, &out)
}

// GetApp calls GET /apps/{appId}.
func (c *Client) GetApp(ctx context.Context, appID int) (*App, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%d", appID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		App App `json:"app"`
	}
	return &out.App, c.do(req, &out)
}

// UpdateApp calls PUT /apps/{appId}.
func (c *Client) UpdateApp(ctx context.Context, appID int, body UpdateAppRequest) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/apps/%d", appID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// DeleteApp calls DELETE /apps/{appId}.
func (c *Client) DeleteApp(ctx context.Context, appID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/apps/%d", appID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// AppAction calls POST /apps/{appId}/actions.
func (c *Client) AppAction(ctx context.Context, appID int, actionType string) (*App, error) {
	body := map[string]string{"type": actionType}
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/apps/%d/actions", appID), body)
	if err != nil {
		return nil, err
	}
	var out struct {
		App App `json:"app"`
	}
	return &out.App, c.do(req, &out)
}

// ---------------------------------------------------------------------------
// AppComponent
// ---------------------------------------------------------------------------

// AppComponent represents a Level27 app component.
type AppComponent struct {
	ID                     int             `json:"id"`
	Name                   string          `json:"name"`
	Category               string          `json:"category"`
	AppComponentType       string          `json:"appcomponenttype"`
	Status                 string          `json:"status"`
	StatusCategory         string          `json:"statusCategory"`
	AppComponentParameters json.RawMessage `json:"appcomponentparameters"`
	App                    *Ref            `json:"app"`
	Attachment             *Ref            `json:"attachment"`
}

// CreateAppComponentRequest is the body for POST /apps/{appId}/components.
type CreateAppComponentRequest struct {
	Name             string `json:"name"`
	AppComponentType string `json:"appcomponenttype"`
	System           int    `json:"system,omitempty"`
	SystemGroup      int    `json:"systemgroup,omitempty"`
	LimitGroup       string `json:"limitGroup,omitempty"`
	// Type-specific parameters
	Version     string `json:"version,omitempty"`
	Path        string `json:"path,omitempty"`
	Pass        string `json:"pass,omitempty"`
	SSHKeys     []int  `json:"sshkeys,omitempty"`
	ExtraConfig string `json:"extraconfig,omitempty"`
	Jailed      *bool  `json:"jailed,omitempty"`
}

// UpdateAppComponentRequest is the body for PUT/PATCH /apps/{appId}/components/{componentId}.
type UpdateAppComponentRequest struct {
	Name        string `json:"name"`
	LimitGroup  string `json:"limitGroup,omitempty"`
	Version     string `json:"version,omitempty"`
	Path        string `json:"path,omitempty"`
	ExtraConfig string `json:"extraconfig,omitempty"`
}

// ListAppComponents calls GET /apps/{appId}/components.
func (c *Client) ListAppComponents(ctx context.Context, appID int) ([]AppComponent, error) {
	q := url.Values{}
	q.Set("limit", "100")
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%d/components?%s", appID, q.Encode()), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Components []AppComponent `json:"components"`
	}
	return out.Components, c.do(req, &out)
}

// CreateAppComponent calls POST /apps/{appId}/components.
func (c *Client) CreateAppComponent(ctx context.Context, appID int, body CreateAppComponentRequest) (*AppComponent, error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/apps/%d/components", appID), body)
	if err != nil {
		return nil, err
	}
	var out struct {
		Component AppComponent `json:"component"`
	}
	return &out.Component, c.do(req, &out)
}

// GetAppComponent calls GET /apps/{appId}/components/{componentId}.
func (c *Client) GetAppComponent(ctx context.Context, appID, componentID int) (*AppComponent, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%d/components/%d", appID, componentID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Component AppComponent `json:"component"`
	}
	return &out.Component, c.do(req, &out)
}

// UpdateAppComponent calls PATCH /apps/{appId}/components/{componentId}.
func (c *Client) UpdateAppComponent(ctx context.Context, appID, componentID int, body UpdateAppComponentRequest) error {
	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/apps/%d/components/%d", appID, componentID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// WaitForComponentStatus polls GET /apps/{appId}/components/{componentId} until
// the status is no longer a transitional state (e.g. "updating", "to_create",
// "creating", "to_delete", "deleting"). Returns the final component or an error.
// A non-green statusCategory ("red") is treated as a failure.
func (c *Client) WaitForComponentStatus(ctx context.Context, appID, componentID int) (*AppComponent, error) {
	transitional := map[string]bool{
		"to_create": true,
		"creating":  true,
		"to_update": true,
		"updating":  true,
		"to_delete": true,
		"deleting":  true,
	}
	for {
		comp, err := c.GetAppComponent(ctx, appID, componentID)
		if err != nil {
			return nil, err
		}
		if !transitional[comp.Status] {
			if comp.StatusCategory == "red" {
				return nil, fmt.Errorf("component reached error status: %s", comp.Status)
			}
			return comp, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}

// DeleteAppComponent calls DELETE /apps/{appId}/components/{componentId}.
func (c *Client) DeleteAppComponent(ctx context.Context, appID, componentID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/apps/%d/components/%d", appID, componentID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---------------------------------------------------------------------------
// AppComponentURL
// ---------------------------------------------------------------------------

// AppComponentURL represents a URL attached to a component.
type AppComponentURL struct {
	ID             int    `json:"id"`
	Content        string `json:"content"`
	HTTPS          bool   `json:"https"`
	Status         string `json:"status"`
	SSLForce       bool   `json:"sslForce"`
	Type           string `json:"type"`
	Caching        bool   `json:"caching"`
	Authentication bool   `json:"authentication"`
}

// CreateAppComponentURLRequest is the body for POST /apps/{appId}/components/{componentId}/urls.
type CreateAppComponentURLRequest struct {
	Content            string `json:"content"`
	SSLForce           bool   `json:"sslForce"`
	HandleDNS          bool   `json:"handleDns"`
	AutoSSLCertificate bool   `json:"autoSslCertificate"`
	Caching            bool   `json:"caching"`
	Authentication     bool   `json:"authentication"`
	SSLCertificate     int    `json:"sslCertificate,omitempty"`
}

// UpdateAppComponentURLRequest is the body for PUT /apps/{appId}/components/{componentId}/urls/{urlId}.
type UpdateAppComponentURLRequest struct {
	Content        string `json:"content"`
	SSLForce       bool   `json:"sslForce"`
	HandleDNS      bool   `json:"handleDns"`
	Caching        bool   `json:"caching"`
	Authentication bool   `json:"authentication"`
	SSLCertificate int    `json:"sslCertificate,omitempty"`
}

// CreateAppComponentURL calls POST /apps/{appId}/components/{componentId}/urls.
func (c *Client) CreateAppComponentURL(ctx context.Context, appID, componentID int, body CreateAppComponentURLRequest) (*AppComponentURL, error) {
	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/apps/%d/components/%d/urls", appID, componentID), body)
	if err != nil {
		return nil, err
	}
	var out struct {
		URL AppComponentURL `json:"url"`
	}
	return &out.URL, c.do(req, &out)
}

// GetAppComponentURL calls GET /apps/{appId}/components/{componentId}/urls/{urlId}.
func (c *Client) GetAppComponentURL(ctx context.Context, appID, componentID, urlID int) (*AppComponentURL, error) {
	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/apps/%d/components/%d/urls/%d", appID, componentID, urlID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		URL AppComponentURL `json:"url"`
	}
	return &out.URL, c.do(req, &out)
}

// UpdateAppComponentURL calls PUT /apps/{appId}/components/{componentId}/urls/{urlId}.
func (c *Client) UpdateAppComponentURL(ctx context.Context, appID, componentID, urlID int, body UpdateAppComponentURLRequest) error {
	req, err := c.newRequest(ctx, http.MethodPut,
		fmt.Sprintf("/apps/%d/components/%d/urls/%d", appID, componentID, urlID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// DeleteAppComponentURL calls DELETE /apps/{appId}/components/{componentId}/urls/{urlId}.
func (c *Client) DeleteAppComponentURL(ctx context.Context, appID, componentID, urlID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/apps/%d/components/%d/urls/%d", appID, componentID, urlID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ---------------------------------------------------------------------------
// SSL Certificate
// ---------------------------------------------------------------------------

// SSLCertificate represents a Level27 SSL certificate.
type SSLCertificate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SSLType   string `json:"sslType"`
	SSLStatus string `json:"sslStatus"`
	Status    string `json:"status"`
}

// CreateSSLCertificateRequest is the body for POST /apps/{appId}/sslcertificates.
type CreateSSLCertificateRequest struct {
	Name                   string `json:"name"`
	SSLType                string `json:"sslType,omitempty"`
	AutoSSLCertificateURLs string `json:"autoSslCertificateUrls,omitempty"`
	SSLKey                 string `json:"sslKey,omitempty"`
	SSLCrt                 string `json:"sslCrt,omitempty"`
	SSLCabundle            string `json:"sslCabundle,omitempty"`
	AutoURLLink            bool   `json:"autoUrlLink"`
	SSLForce               bool   `json:"sslForce"`
}

// CreateSSLCertificate calls POST /apps/{appId}/sslcertificates.
func (c *Client) CreateSSLCertificate(ctx context.Context, appID int, body CreateSSLCertificateRequest) (*SSLCertificate, error) {
	req, err := c.newRequest(ctx, http.MethodPost,
		fmt.Sprintf("/apps/%d/sslcertificates", appID), body)
	if err != nil {
		return nil, err
	}
	var out struct {
		SSLCertificate SSLCertificate `json:"sslCertificate"`
	}
	return &out.SSLCertificate, c.do(req, &out)
}

// GetSSLCertificate calls GET /apps/{appId}/sslcertificates/{sslcertificateId}.
func (c *Client) GetSSLCertificate(ctx context.Context, appID, certID int) (*SSLCertificate, error) {
	req, err := c.newRequest(ctx, http.MethodGet,
		fmt.Sprintf("/apps/%d/sslcertificates/%d", appID, certID), nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		SSLCertificate SSLCertificate `json:"sslCertificate"`
	}
	return &out.SSLCertificate, c.do(req, &out)
}

// UpdateSSLCertificate calls PUT /apps/{appId}/sslcertificates/{sslcertificateId}.
func (c *Client) UpdateSSLCertificate(ctx context.Context, appID, certID int, body map[string]string) error {
	req, err := c.newRequest(ctx, http.MethodPut,
		fmt.Sprintf("/apps/%d/sslcertificates/%d", appID, certID), body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// DeleteSSLCertificate calls DELETE /apps/{appId}/sslcertificates/{sslcertificateId}.
func (c *Client) DeleteSSLCertificate(ctx context.Context, appID, certID int) error {
	req, err := c.newRequest(ctx, http.MethodDelete,
		fmt.Sprintf("/apps/%d/sslcertificates/%d", appID, certID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
