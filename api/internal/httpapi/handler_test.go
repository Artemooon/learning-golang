package httpapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artem/vintage-jazz-records-api/internal/httpapi"
	"github.com/artem/vintage-jazz-records-api/internal/records"
)

func TestRecordsAPI(t *testing.T) {
	api := httpapi.NewHandler(records.NewStore(nil))
	body := []byte(`{"id":"a-love-supreme","title":"A Love Supreme","artist":"John Coltrane","year":1965,"label":"Impulse!"}`)

	create := httptest.NewRequest(http.MethodPost, "/records", bytes.NewReader(body))
	created := httptest.NewRecorder()
	api.ServeHTTP(created, create)
	if created.Code != http.StatusCreated {
		t.Fatalf("POST /records status = %d, want %d", created.Code, http.StatusCreated)
	}

	get := httptest.NewRequest(http.MethodGet, "/records/a-love-supreme", nil)
	found := httptest.NewRecorder()
	api.ServeHTTP(found, get)
	if found.Code != http.StatusOK {
		t.Fatalf("GET /records/{id} status = %d, want %d", found.Code, http.StatusOK)
	}

	list := httptest.NewRequest(http.MethodGet, "/records", nil)
	listed := httptest.NewRecorder()
	api.ServeHTTP(listed, list)
	if listed.Code != http.StatusOK {
		t.Fatalf("GET /records status = %d, want %d", listed.Code, http.StatusOK)
	}
}
