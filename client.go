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
	LNProxyError                = errors.New("")
	PaymentHashMismatch         = errors.New("payment hash does not match")
	CustomDescriptionMismatch   = errors.New("custom description does match")
	CustomRoutingBudgetMismatch = errors.New("routing budget not respected")
	InvalidProxyInvoice         = errors.New("invalid proxy invoice")
)

type LNProxy struct {
	url.URL
	http.Client
	BaseMsat uint64
	Ppm      uint64
}

func (x *LNProxy) RequestProxy(invoice string, routing_msat uint64) (proxy_invoice string, err error) {
	params, _ := json.Marshal(struct {
		Invoice     string `json:"invoice"`
		RoutingMsat uint64 `json:"routing_msat,string"`
	}{
		Invoice:     invoice,
		RoutingMsat: routing_msat,
	})
	buf := bytes.NewBuffer(params)
	req, err := http.NewRequest("POST", x.URL.String(), buf)
	if err != nil {
		return "", err
	}
	resp, err := x.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		r := struct {
			Reason string `json:"reason"`
		}{}
		err = dec.Decode(&r)
		if err != nil && err != io.EOF {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("malformed lnproxy response: %s", string(body))
		}
		return "", errors.Join(LNProxyError, errors.New(r.Reason))

	}
	r := struct {
		ProxyInvoice string `json:"proxy_invoice"`
	}{}
	err = dec.Decode(&r)
	if err != nil && err != io.EOF {
		return "", err
	}
	return r.ProxyInvoice, nil
}

func ValidateProxyInvoice(invoice, proxy_invoice string, routing_msat uint64) (bool, error) {
	original, err := ParseInvoice([]byte(invoice))
	if err != nil {
		return false, errors.New("invalid original invoice")
	}
	proxy, err := ParseInvoice([]byte(proxy_invoice))
	if err != nil {
		return false, InvalidProxyInvoice
	}
	if bytes.Compare(original.PaymentHash, proxy.PaymentHash) != 0 {
		return false, PaymentHashMismatch
	}
	if original.DescriptionHash != proxy.DescriptionHash {
		return false, CustomDescriptionMismatch
	}
	if bytes.Compare(original.Description, proxy.Description) != 0 {
		return false, CustomDescriptionMismatch
	}
	if (original.AmountMsat + routing_msat) != proxy.AmountMsat {
		return false, CustomRoutingBudgetMismatch
	}
	return true, nil
}
