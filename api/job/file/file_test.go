package jobfile

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

func newTestJobFile(t *testing.T, handler http.HandlerFunc) JobFile {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := smclient.NewClient(server.Client(), "user", "secret")
	client.BaseURL = server.URL
	client.Credentials.AccessToken = &smclient.Token{
		Value:          "fake-test-token",
		ExpirationTime: time.Now().Add(1 * time.Hour),
	}

	return NewJobFile(client)
}

func TestAdd_PostsFileToAddPath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/file/add") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got AddRequest
		_ = json.Unmarshal(body, &got)
		if got.FileURI != "a.json" {
			t.Errorf("fileUri = %q, want a.json", got.FileURI)
		}
		if len(got.TargetLocaleIDs) != 1 || got.TargetLocaleIDs[0] != "fr-FR" {
			t.Errorf("targetLocaleIds = %v, want [fr-FR]", got.TargetLocaleIDs)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"successCount":1,"failCount":0}}}`))
	}

	jf := newTestJobFile(t, handler)

	res, err := jf.Add(context.Background(), "p1", "job-1", AddRequest{FileURI: "a.json", TargetLocaleIDs: []string{"fr-FR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SuccessCount != 1 {
		t.Errorf("successCount = %d, want 1", res.SuccessCount)
	}
}

func TestRemove_PostsFileToRemovePath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/file/remove") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got RemoveRequest
		_ = json.Unmarshal(body, &got)
		if got.FileURI != "a.json" {
			t.Errorf("fileUri = %q, want a.json", got.FileURI)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"successCount":1,"failCount":0}}}`))
	}

	jf := newTestJobFile(t, handler)

	if _, err := jf.Remove(context.Background(), "p1", "job-1", RemoveRequest{FileURI: "a.json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestList_ForwardsPagingAndMapsItems(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/files") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("limit") != "50" || q.Get("offset") != "100" {
			t.Errorf("limit/offset = %q/%q, want 50/100", q.Get("limit"), q.Get("offset"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{
			"totalCount": 2,
			"items": [
				{"uri":"/a.json","localeIds":["fr-FR"]},
				{"uri":"/b.xml"}
			]
		}}}`))
	}

	jf := newTestJobFile(t, handler)

	page, err := jf.List(context.Background(), "p1", "job-1", 50, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.TotalCount != 2 || len(page.Items) != 2 {
		t.Fatalf("unexpected page: %+v", page)
	}
	if page.Items[0].FileURI != "/a.json" || len(page.Items[0].LocaleIDs) != 1 {
		t.Errorf("Items[0] = %+v, want /a.json [fr-FR]", page.Items[0])
	}
	if page.Items[1].FileURI != "/b.xml" || len(page.Items[1].LocaleIDs) != 0 {
		t.Errorf("Items[1] = %+v, want /b.xml []", page.Items[1])
	}
}

func TestList_NotFoundMapsToErrNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	jf := newTestJobFile(t, handler)

	if _, err := jf.List(context.Background(), "p1", "missing", 500, 0); !errors.Is(err, jobapi.ErrNotFound) {
		t.Errorf("err = %v, want jobapi.ErrNotFound", err)
	}
}

func TestAddRemove_NotFoundMapsToErrNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	jf := newTestJobFile(t, handler)

	if _, err := jf.Add(context.Background(), "p1", "missing", AddRequest{FileURI: "a.json"}); !errors.Is(err, jobapi.ErrNotFound) {
		t.Errorf("Add err = %v, want jobapi.ErrNotFound", err)
	}
	if _, err := jf.Remove(context.Background(), "p1", "missing", RemoveRequest{FileURI: "a.json"}); !errors.Is(err, jobapi.ErrNotFound) {
		t.Errorf("Remove err = %v, want jobapi.ErrNotFound", err)
	}
}

func TestAddRemove_EmptyParamsRejectedBeforeRequest(t *testing.T) {
	called := false
	jf := newTestJobFile(t, func(http.ResponseWriter, *http.Request) { called = true })

	cases := []struct {
		name                       string
		projectID, jobUID, fileURI string
	}{
		{"empty project", "", "job-1", "a.json"},
		{"empty job", "p1", "", "a.json"},
		{"empty fileUri", "p1", "job-1", ""},
	}
	for _, c := range cases {
		t.Run("add/"+c.name, func(t *testing.T) {
			_, err := jf.Add(context.Background(), c.projectID, c.jobUID, AddRequest{FileURI: c.fileURI})
			if !smerror.IsErrEmptyParam(err) {
				t.Errorf("Add err = %v, want empty-param error", err)
			}
		})
		t.Run("remove/"+c.name, func(t *testing.T) {
			_, err := jf.Remove(context.Background(), c.projectID, c.jobUID, RemoveRequest{FileURI: c.fileURI})
			if !smerror.IsErrEmptyParam(err) {
				t.Errorf("Remove err = %v, want empty-param error", err)
			}
		})
	}

	if called {
		t.Error("server was called despite empty-param validation")
	}
}
