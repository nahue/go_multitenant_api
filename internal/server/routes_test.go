package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"test_go_api/internal/database"
	"test_go_api/internal/testutil"
)

func TestMain(m *testing.M) {
	os.Exit(testutil.SetupTestDB(m))
}

func TestHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.HelloWorldHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestTenantNotProvidedInSubdomain(t *testing.T) {
	s := &Server{
		db: database.New(),
	}
	handler := s.RegisterRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test without tenant in subdomain
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", resp.Status)
	}
}

func TestTenantProvidedInSubdomain(t *testing.T) {
	s := &Server{
		db: database.New(),
	}
	handler := s.RegisterRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	serverURL := strings.Replace(server.URL, "127.0.0.1", "test.app.localhost", 1)
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}
