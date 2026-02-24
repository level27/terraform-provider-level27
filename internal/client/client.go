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

// SystemCookbook represents a single cookbook installed on a system.
type SystemCookbook struct {
	ID                 int                          `json:"id"`
	CookbookType       string                       `json:"cookbooktype"`
	CookbookParameters map[string]CookbookParameter `json:"cookbookparameters"`
}

// CookbookParameter holds the value for a single cookbook parameter.
// The value can be a string, number, bool, or array of strings.
type CookbookParameter struct {
	Value json.RawMessage `json:"value"`
}

// Versions returns the list of versions for this cookbook parameter.
// The API returns either a string or []string.
func (cp CookbookParameter) Versions() []string {
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

// System represents a Level27 system.
type System struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Cookbooks []SystemCookbook `json:"cookbooks"`
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

// FindCookbookVersion looks up the highest available version for the given
// cookbook type on this system. Returns "" if not found.
// The API stores versions under either "version" or "versions" key.
func (s *System) FindCookbookVersion(cookbookType string) string {
	for _, cb := range s.Cookbooks {
		if cb.CookbookType != cookbookType {
			continue
		}
		var all []string
		for _, key := range []string{"version", "versions"} {
			if vp, ok := cb.CookbookParameters[key]; ok {
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
