package client

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
)

var charSet = []byte("qpzry9x8gf2tvdw0s3jn54khce6mua7l")

var isBech32 = regexp.MustCompile("^lnbc(?:[0-9]+[pnum])?1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]+$")

type InvoiceParts struct {
	AmountMsat      uint64
	PaymentHash     []byte
	Description     []byte
	DescriptionHash bool
	Signature       []byte
}

func ParseInvoice(invoice []byte) (*InvoiceParts, error) {
	invoice = bytes.ToLower(invoice)
	pos := bytes.LastIndexByte(invoice, byte('1'))
	if pos == -1 || !isBech32.Match(invoice) {
		return nil, errors.New("invalid invoice")
	}

	var invoice_parts InvoiceParts
	var err error
	if pos > 4 {
		invoice_parts.AmountMsat, err = strconv.ParseUint(string(invoice[4:pos-1]), 10, 64)
		if err != nil {
			return nil, err
		}
		switch invoice[pos-1] {
		case byte('p'):
			invoice_parts.AmountMsat /= 10
		case byte('n'):
			invoice_parts.AmountMsat *= 100
		case byte('u'):
			invoice_parts.AmountMsat *= 100_000
		case byte('m'):
			invoice_parts.AmountMsat *= 100_000_000
		}
	}
	for i := pos + 8; i < len(invoice); {
		data_length := bytes.Index(charSet, invoice[i+1:i+2])*32 + bytes.Index(charSet, invoice[i+2:i+3])
		if invoice[i] == byte('p') {
			invoice_parts.PaymentHash = invoice[i+3 : i+3+data_length]
		}
		if invoice[i] == byte('d') {
			invoice_parts.DescriptionHash = false
			invoice_parts.Description = invoice[i+3 : i+3+data_length]
		}
		if invoice[i] == byte('h') {
			invoice_parts.DescriptionHash = true
			invoice_parts.Description = invoice[i+3 : i+3+data_length]
		}
		i += 3 + data_length
	}
	invoice_parts.Signature = invoice[len(invoice)-110 : len(invoice)-6]
	return &invoice_parts, nil
}
