package cli

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type service struct {
	addr string
	auth string
	http *http.Client
}

func newService(addr, auth string) *service {
	return &service{
		addr: addr,
		auth: auth,
		http: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				DisableKeepAlives:     true,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   -1,
				IdleConnTimeout:       90 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (c *service) newRequest(method, path string, form interface{}) (*http.Response, error) {
	var body *bytes.Reader
	b, err := json.Marshal(form)
	if err == nil {
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.addr+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	c.maybeAuthorize(req)
	return c.http.Do(req)
}

func (c *service) maybeAuthorize(req *http.Request) {
	if c.auth != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth)
	}
}

func (c *service) parseRequest(method, path string, form, view interface{}) error {
	resp, err := c.newRequest(method, path, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if view == nil {
		return nil
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return json.NewDecoder(resp.Body).Decode(&view)
	case http.StatusNoContent:
		return nil
	}
	var verr errorResponse
	err = json.NewDecoder(resp.Body).Decode(&verr)
	if err != nil {
		return err
	}
	return verr
}

func (c *service) getMachine(id string) (*getMachineResponse, error) {
	var view *getMachineResponse
	return view, c.parseRequest(http.MethodGet, "/machines/"+id, nil, &view)
}

func (c *service) initMachine(form flagsInit) (string, error) {
	var view getMachineResponse
	err := c.parseRequest(http.MethodPost, "/machines", form, &view)
	return view.ID, err
}

func (c *service) cancelMachine(id string) error {
	return c.parseRequest(http.MethodPost, "/machines/"+id+"/cancel", nil, nil)
}

func (c *service) renewMachine(id string) error {
	return c.parseRequest(http.MethodPost, "/machines/"+id+"/renew", nil, nil)
}

func (c *service) destroyMachine(id string, form flagsDestroy) error {
	return c.parseRequest(http.MethodDelete, "/machines/"+id, form, nil)
}

type getMachineResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	IPv4      string    `json:"ipv4"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}

type errorResponse struct {
	Code      int    `json:"code"`
	Title     string `json:"title"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id"`
}

func (e errorResponse) Error() string {
	return e.Message
}
