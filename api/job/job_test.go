package job

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

func newTestJob(t *testing.T, handler http.HandlerFunc) (Job, *httptest.Server) {
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

	return NewJob(client), server
}

func TestGetJob_PopulatesTargetLocaleIDs(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/jobs/jobUid-123") {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"code": "SUCCESS",
				"data": {
					"jobName": "My Job",
					"translationJobUid": "jobUid-123",
					"targetLocaleIds": ["fr-FR", "de-DE", "es-ES"]
				}
			}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.GetJob(context.Background(), "projectId-xyz", "jobUid-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"fr-FR", "de-DE", "es-ES"}
	if !equalStringSlices(resp.TargetLocaleIDs, want) {
		t.Errorf("TargetLocaleIDs: got %v, want %v", resp.TargetLocaleIDs, want)
	}
	if resp.JobName != "My Job" {
		t.Errorf("JobName: got %q, want %q", resp.JobName, "My Job")
	}
	if resp.TranslationJobUID != "jobUid-123" {
		t.Errorf("TranslationJobUID: got %q, want %q", resp.TranslationJobUID, "jobUid-123")
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestListFiles_ReturnsPageAndForwardsLimitOffset(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/jobs/jobUid-123/files") {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("offset"); got != "100" {
			t.Errorf("offset: got %q, want %q", got, "100")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("limit: got %q, want %q", got, "50")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"code": "SUCCESS",
				"data": {
					"totalCount": 750,
					"items": [
						{"uri": "/path/a.json", "fileUid": "abc123", "localeIds": ["fr-FR"]},
						{"uri": "/path/b.xml", "fileUid": "def456"}
					]
				}
			}
		}`))
	}

	j, _ := newTestJob(t, handler)

	page, err := j.ListFiles(context.Background(), "projectId-xyz", "jobUid-123", 50, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.TotalCount != 750 {
		t.Errorf("TotalCount = %d, want 750", page.TotalCount)
	}
	if len(page.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(page.Items))
	}
	if page.Items[0].FileURI != "/path/a.json" {
		t.Errorf("Items[0].FileURI = %q, want %q", page.Items[0].FileURI, "/path/a.json")
	}
	if !equalStringSlices(page.Items[0].LocaleIDs, []string{"fr-FR"}) {
		t.Errorf("Items[0].LocaleIDs = %v, want [fr-FR]", page.Items[0].LocaleIDs)
	}
	if page.Items[1].FileURI != "/path/b.xml" {
		t.Errorf("Items[1].FileURI = %q, want %q", page.Items[1].FileURI, "/path/b.xml")
	}
	if len(page.Items[1].LocaleIDs) != 0 {
		t.Errorf("Items[1].LocaleIDs = %v, want empty", page.Items[1].LocaleIDs)
	}
}

func TestListFiles_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	j, _ := newTestJob(t, handler)

	_, err := j.ListFiles(context.Background(), "p", "missing-job", 500, 0)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestListFiles_Empty(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"totalCount":0,"items":[]}}}`))
	}

	j, _ := newTestJob(t, handler)

	page, err := j.ListFiles(context.Background(), "p", "j", 500, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Items) != 0 {
		t.Errorf("Items = %v, want empty", page.Items)
	}
	if page.TotalCount != 0 {
		t.Errorf("TotalCount = %d, want 0", page.TotalCount)
	}
}

func TestListFiles_ServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"response":{"code":"GENERAL_ERROR"}}`))
	}

	j, _ := newTestJob(t, handler)

	_, err := j.ListFiles(context.Background(), "p", "j", 500, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrNotFound) {
		t.Errorf("err = ErrNotFound, want generic error")
	}
	if !strings.Contains(err.Error(), "failed to list job files") {
		t.Errorf("err = %v, want wrapped %q", err, "failed to list job files")
	}
}
