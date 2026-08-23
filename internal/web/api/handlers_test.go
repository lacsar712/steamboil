package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lacsar712/steamboil/internal/app"
	"github.com/lacsar712/steamboil/internal/clock"
	"github.com/lacsar712/steamboil/internal/config"
	"github.com/lacsar712/steamboil/internal/web/api"
)

func testServer(t *testing.T) *api.Server {
	t.Helper()
	clk := clock.NewManual(clock.NewReal().Now())
	cfg := config.Default("API-1")
	a, err := app.BootstrapWithClock(cfg, clk)
	if err != nil {
		t.Fatal(err)
	}
	return api.NewServer(a)
}

func TestHealthEndpoint(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var body struct {
		OK   bool `json:"ok"`
		Data map[string]string
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.OK || body.Data["status"] != "ok" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestSnapshotEndpoint(t *testing.T) {
	s := testServer(t)
	req := httptest.NewRequest(http.MethodGet, "/snapshot", nil)
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", rec.Code, rec.Body.String())
	}
}
