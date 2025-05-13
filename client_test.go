package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestProxy(t *testing.T) {
	// Setup a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Parse and validate the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		var reqBody struct {
			Invoice     string `json:"invoice"`
			RoutingMsat uint64 `json:"routing_msat,string"`
		}

		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		if reqBody.Invoice != "test-invoice" {
			t.Errorf("Expected invoice 'test-invoice', got '%s'", reqBody.Invoice)
		}

		if reqBody.RoutingMsat != 1500 {
			t.Errorf("Expected routing_msat 1500, got %d", reqBody.RoutingMsat)
		}

		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := struct {
			ProxyInvoice string `json:"proxy_invoice"`
		}{
			ProxyInvoice: "proxy-test-invoice",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Setup a test buffer to capture logs
	var logBuffer bytes.Buffer
	testLogger := NewLogger(LevelDebug, &logBuffer)

	// Create a client with the test server URL
	serverURL, _ := url.Parse(server.URL)
	client := NewLNProxy(*serverURL, 1000, 500).WithLogger(testLogger)

	// Call the method under test
	proxyInvoice, err := client.RequestProxy("test-invoice", 1500)
	if err != nil {
		t.Fatalf("RequestProxy failed: %v", err)
	}

	// Check the result
	if proxyInvoice != "proxy-test-invoice" {
		t.Errorf("Expected 'proxy-test-invoice', got '%s'", proxyInvoice)
	}

	// Check that logs were generated
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "Requesting proxy invoice") {
		t.Error("Expected log message 'Requesting proxy invoice' not found")
	}
	if !strings.Contains(logOutput, "Successfully received proxy invoice") {
		t.Error("Expected log message 'Successfully received proxy invoice' not found")
	}
}

func TestRequestProxyError(t *testing.T) {
	// Setup a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send an error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := struct {
			Reason string `json:"reason"`
		}{
			Reason: "Invalid invoice",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Setup a test buffer to capture logs
	var logBuffer bytes.Buffer
	testLogger := NewLogger(LevelDebug, &logBuffer)

	// Create a client with the test server URL
	serverURL, _ := url.Parse(server.URL)
	client := NewLNProxy(*serverURL, 1000, 500).WithLogger(testLogger)

	// Call the method under test
	_, err := client.RequestProxy("invalid-invoice", 1500)
	
	// Check that we got an error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check the error message
	if !strings.Contains(err.Error(), "Invalid invoice") {
		t.Errorf("Expected error containing 'Invalid invoice', got '%s'", err.Error())
	}

	// Check that error was logged
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "LNProxy error: Invalid invoice") {
		t.Error("Expected log message 'LNProxy error: Invalid invoice' not found")
	}
}
