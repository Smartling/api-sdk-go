package job

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestListFiles_SinglePage(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/jobs/jobUid-123/files") {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("offset"); got != "0" {
			t.Errorf("offset on first request: got %q, want %q", got, "0")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"code": "SUCCESS",
				"data": {
					"totalCount": 2,
					"items": [
						{"uri": "/path/a.json", "fileUid": "abc123", "localeIds": ["fr-FR"]},
						{"uri": "/path/b.xml", "fileUid": "def456"}
					]
				}
			}
		}`))
	}

	j, _ := newTestJob(t, handler)

	files, err := j.ListFiles(context.Background(), "projectId-xyz", "jobUid-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("len(files) = %d, want 2", len(files))
	}
	if files[0].FileURI != "/path/a.json" {
		t.Errorf("files[0].FileURI = %q, want %q", files[0].FileURI, "/path/a.json")
	}
	if !equalStringSlices(files[0].LocaleIDs, []string{"fr-FR"}) {
		t.Errorf("files[0].LocaleIDs = %v, want [fr-FR]", files[0].LocaleIDs)
	}
	if files[1].FileURI != "/path/b.xml" {
		t.Errorf("files[1].FileURI = %q, want %q", files[1].FileURI, "/path/b.xml")
	}
	if len(files[1].LocaleIDs) != 0 {
		t.Errorf("files[1].LocaleIDs = %v, want empty", files[1].LocaleIDs)
	}
}

func TestListFiles_MultiPage(t *testing.T) {
	const total = 750

	var pagesServed int
	handler := func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit != 500 {
			t.Errorf("limit: got %d, want 500", limit)
		}

		pagesServed++

		count := total - offset
		if count > limit {
			count = limit
		}

		var items strings.Builder
		for i := 0; i < count; i++ {
			if i > 0 {
				items.WriteString(",")
			}
			items.WriteString(`{"uri":"/f`)
			items.WriteString(strconv.Itoa(offset + i))
			items.WriteString(`.json"}`)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"totalCount":` +
			strconv.Itoa(total) + `,"items":[` + items.String() + `]}}}`))
	}

	j, _ := newTestJob(t, handler)

	files, err := j.ListFiles(context.Background(), "p", "j")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != total {
		t.Fatalf("len(files) = %d, want %d", len(files), total)
	}
	if pagesServed != 2 {
		t.Errorf("pagesServed = %d, want 2 (500 + 250 split)", pagesServed)
	}
	if files[0].FileURI != "/f0.json" {
		t.Errorf("files[0] = %q, want /f0.json", files[0].FileURI)
	}
	if files[total-1].FileURI != "/f749.json" {
		t.Errorf("files[last] = %q, want /f749.json", files[total-1].FileURI)
	}
}

func TestListFiles_NotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	j, _ := newTestJob(t, handler)

	_, err := j.ListFiles(context.Background(), "p", "missing-job")
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

	files, err := j.ListFiles(context.Background(), "p", "j")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if files != nil {
		t.Errorf("files = %v, want nil", files)
	}
}

func TestListFiles_ServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"response":{"code":"GENERAL_ERROR"}}`))
	}

	j, _ := newTestJob(t, handler)

	_, err := j.ListFiles(context.Background(), "p", "j")
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
