package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	LNProxyError                = errors.New("lnproxy error")
	PaymentHashMismatch         = errors.New("payment hash does not match")
	DescriptionMismatch         = errors.New("description does match")
	CustomRoutingBudgetMismatch = errors.New("routing budget not respected")
	DestinationNotProxied       = errors.New("destination is not obscured")
	InvalidProxyInvoice         = errors.New("invalid proxy invoice")
)

type LNProxy struct {
	url.URL
	http.Client
	BaseMsat uint64
	Ppm      uint64
	logger   *Logger
}

// NewLNProxy creates a new LNProxy client with the default logger
func NewLNProxy(baseURL url.URL, baseMsat, ppm uint64) *LNProxy {
	return &LNProxy{
		URL:      baseURL,
		Client:   http.Client{},
		BaseMsat: baseMsat,
		Ppm:      ppm,
		logger:   DefaultLogger().WithComponent("LNProxy"),
	}
}

// WithLogger sets a custom logger for the LNProxy client
func (x *LNProxy) WithLogger(logger *Logger) *LNProxy {
	x.logger = logger.WithComponent("LNProxy")
	return x
}

func (x *LNProxy) RequestProxy(invoice string, routing_msat uint64) (proxy_invoice string, err error) {
	x.logger.Debug("Requesting proxy invoice for %s with routing budget %d msat", invoice, routing_msat)
	
	params, _ := json.Marshal(struct {
		Invoice     string `json:"invoice"`
		RoutingMsat string `json:"routing_msat"`
	}{
		Invoice:     invoice,
		RoutingMsat: fmt.Sprintf("%d", routing_msat),
	})
	
	buf := bytes.NewBuffer(params)
	req, err := http.NewRequest("POST", x.URL.String(), buf)
	if err != nil {
		x.logger.Error("Failed to create HTTP request: %v", err)
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	x.logger.Debug("Sending request to %s", x.URL.String())
	resp, err := x.Client.Do(req)
	if err != nil {
		x.logger.Error("HTTP request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	
	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		x.logger.Warn("Received non-OK status code: %d", resp.StatusCode)
		r := struct {
			Reason string `json:"reason"`
			Status string `json:"status,omitempty"`
		}{}
		err = dec.Decode(&r)
		if err != nil && err != io.EOF {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				x.logger.Error("Failed to read response body: %v", err)
				return "", err
			}
			x.logger.Error("Malformed lnproxy response: %s", string(body))
			return "", fmt.Errorf("malformed lnproxy response: %s", string(body))
		}
		x.logger.Error("LNProxy error: %s", r.Reason)
		return "", errors.Join(LNProxyError, errors.New(r.Reason))
	}
	
	r := struct {
		ProxyInvoice string `json:"proxy_invoice"`
	}{}
	err = dec.Decode(&r)
	if err != nil && err != io.EOF {
		x.logger.Error("Failed to decode successful response: %v", err)
		return "", err
	}
	
	x.logger.Debug("Successfully received proxy invoice: %s", r.ProxyInvoice)
	return r.ProxyInvoice, nil
}

func ValidateProxyInvoice(invoice, proxy_invoice string, routing_msat uint64) (bool, error) {
	logger := DefaultLogger().WithComponent("Validator")
	logger.Debug("Validating proxy invoice against original invoice")
	
	original, err := ParseInvoice([]byte(invoice))
	if err != nil {
		logger.Error("Failed to parse original invoice: %v", err)
		return false, errors.New("invalid original invoice")
	}
	
	proxy, err := ParseInvoice([]byte(proxy_invoice))
	if err != nil {
		logger.Error("Failed to parse proxy invoice: %v", err)
		return false, InvalidProxyInvoice
	}
	
	logger.Debug("Original amount: %d msat, Proxy amount: %d msat", original.AmountMsat, proxy.AmountMsat)
	
	if bytes.Compare(original.PaymentHash, proxy.PaymentHash) != 0 {
		logger.Error("Payment hash mismatch")
		return false, PaymentHashMismatch
	}
	
	if original.DescriptionHash != proxy.DescriptionHash {
		logger.Error("Description hash mismatch")
		return false, DescriptionMismatch
	}
	
	if bytes.Compare(original.Description, proxy.Description) != 0 {
		logger.Error("Description mismatch")
		return false, DescriptionMismatch
	}
	
	if (original.AmountMsat + routing_msat) != proxy.AmountMsat {
		logger.Error("Routing budget mismatch: expected %d, got %d", 
			original.AmountMsat+routing_msat, proxy.AmountMsat)
		return false, CustomRoutingBudgetMismatch
	}
	
	if bytes.Compare(original.Signature, proxy.Signature) == 0 {
		logger.Error("Destination not proxied (signatures match)")
		return false, DestinationNotProxied
	}
	
	logger.Debug("Proxy invoice validation successful")
	return true, nil
}
