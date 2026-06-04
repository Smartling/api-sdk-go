package locale

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

func newTestJobLocale(t *testing.T, handler http.HandlerFunc) JobLocale {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := smclient.NewClient(server.Client(), "user", "secret")
	client.BaseURL = server.URL
	// Pre-populate a non-expiring access token so requests skip the live
	// /auth-api/v2/authenticate handshake.
	client.Credentials.AccessToken = &smclient.Token{
		Value:          "fake-test-token",
		ExpirationTime: time.Now().Add(1 * time.Hour),
	}

	return NewJobLocale(client)
}

func okEnvelope(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":null}}`))
}

func TestAdd_PostsToLocalePath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/p1/jobs/job-1/locales/fr-FR") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		okEnvelope(w)
	}

	jl := newTestJobLocale(t, handler)

	if err := jl.Add(context.Background(), "p1", "job-1", "fr-FR"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemove_DeletesLocalePath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/p1/jobs/job-1/locales/fr-FR") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		okEnvelope(w)
	}

	jl := newTestJobLocale(t, handler)

	if err := jl.Remove(context.Background(), "p1", "job-1", "fr-FR"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddRemove_EmptyParamsRejectedBeforeRequest(t *testing.T) {
	called := false
	jl := newTestJobLocale(t, func(http.ResponseWriter, *http.Request) { called = true })

	cases := []struct {
		name                            string
		projectID, jobUID, targetLocale string
	}{
		{"empty project", "", "job-1", "fr-FR"},
		{"empty job", "p1", "", "fr-FR"},
		{"empty locale", "p1", "job-1", ""},
	}
	for _, c := range cases {
		t.Run("add/"+c.name, func(t *testing.T) {
			err := jl.Add(context.Background(), c.projectID, c.jobUID, c.targetLocale)
			if !smerror.IsErrEmptyParam(err) {
				t.Errorf("Add err = %v, want empty-param error", err)
			}
		})
		t.Run("remove/"+c.name, func(t *testing.T) {
			err := jl.Remove(context.Background(), c.projectID, c.jobUID, c.targetLocale)
			if !smerror.IsErrEmptyParam(err) {
				t.Errorf("Remove err = %v, want empty-param error", err)
			}
		})
	}

	if called {
		t.Error("server was called despite empty-param validation")
	}
}

func TestAdd_NotFoundIsWrapped(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	jl := newTestJobLocale(t, handler)

	err := jl.Add(context.Background(), "p1", "job-1", "fr-FR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to add locale to job") {
		t.Errorf("err = %v, want wrapped %q", err, "failed to add locale to job")
	}
}

func TestRemove_ServerErrorIsWrapped(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"response":{"code":"GENERAL_ERROR"}}`))
	}

	jl := newTestJobLocale(t, handler)

	err := jl.Remove(context.Background(), "p1", "job-1", "fr-FR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to remove locale from job") {
		t.Errorf("err = %v, want wrapped %q", err, "failed to remove locale from job")
	}
}
