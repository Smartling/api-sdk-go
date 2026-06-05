package jobfile

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestAddRemove_EmptyParamsRejectedBeforeRequest(t *testing.T) {
	called := false
	jf := newTestJobFile(t, func(http.ResponseWriter, *http.Request) { called = true })

	if _, err := jf.Add(context.Background(), "p1", "job-1", AddRequest{}); !smerror.IsErrEmptyParam(err) {
		t.Errorf("Add with no fileUri err = %v, want empty-param error", err)
	}
	if _, err := jf.Remove(context.Background(), "", "job-1", RemoveRequest{FileURI: "a.json"}); !smerror.IsErrEmptyParam(err) {
		t.Errorf("Remove with empty project err = %v, want empty-param error", err)
	}

	if called {
		t.Error("server was called despite empty-param validation")
	}
}
