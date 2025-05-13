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
	logger := DefaultLogger().WithComponent("InvoiceParser")
	logger.Debug("Parsing invoice: %s", string(invoice[:min(len(invoice), 40)]) + "...")
	
	invoice = bytes.ToLower(invoice)
	pos := bytes.LastIndexByte(invoice, byte('1'))
	if pos == -1 || !isBech32.Match(invoice) {
		logger.Error("Invalid invoice format")
		return nil, errors.New("invalid invoice")
	}

	var invoice_parts InvoiceParts
	var err error
	if pos > 4 {
		amountStr := string(invoice[4 : pos-1])
		logger.Debug("Parsing amount from: %s", amountStr)
		invoice_parts.AmountMsat, err = strconv.ParseUint(amountStr, 10, 64)
		if err != nil {
			logger.Error("Failed to parse amount: %v", err)
			return nil, err
		}
		
		unit := invoice[pos-1]
		logger.Debug("Amount unit: %c", unit)
		switch unit {
		case byte('p'):
			invoice_parts.AmountMsat /= 10
		case byte('n'):
			invoice_parts.AmountMsat *= 100
		case byte('u'):
			invoice_parts.AmountMsat *= 100_000
		case byte('m'):
			invoice_parts.AmountMsat *= 100_000_000
		}
		logger.Debug("Calculated amount: %d msat", invoice_parts.AmountMsat)
	}
	
	logger.Debug("Parsing invoice data fields")
	for i := pos + 8; i < len(invoice); {
		data_length := bytes.Index(charSet, invoice[i+1:i+2])*32 + bytes.Index(charSet, invoice[i+2:i+3])
		logger.Debug("Found field type %c with length %d", invoice[i], data_length)
		
		if invoice[i] == byte('p') {
			invoice_parts.PaymentHash = invoice[i+3 : i+3+data_length]
			logger.Debug("Payment hash found (length: %d)", len(invoice_parts.PaymentHash))
		}
		if invoice[i] == byte('d') {
			invoice_parts.DescriptionHash = false
			invoice_parts.Description = invoice[i+3 : i+3+data_length]
			logger.Debug("Description found: %s", string(invoice_parts.Description))
		}
		if invoice[i] == byte('h') {
			invoice_parts.DescriptionHash = true
			invoice_parts.Description = invoice[i+3 : i+3+data_length]
			logger.Debug("Description hash found (length: %d)", len(invoice_parts.Description))
		}
		i += 3 + data_length
	}
	
	invoice_parts.Signature = invoice[len(invoice)-110 : len(invoice)-6]
	logger.Debug("Signature extracted (length: %d)", len(invoice_parts.Signature))
	
	logger.Debug("Invoice parsing complete")
	return &invoice_parts, nil
}

// Helper function for min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
