package jobstring

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

func TestAddRemove_NotFoundMapsToErrNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":{"code":"NOT_FOUND","data":null}}`))
	}

	js := newTestJobString(t, handler)

	if _, err := js.Add(context.Background(), "p1", "missing", AddRequest{Hashcodes: []string{"h1"}}); !errors.Is(err, jobapi.ErrNotFound) {
		t.Errorf("Add err = %v, want jobapi.ErrNotFound", err)
	}
	if _, err := js.Remove(context.Background(), "p1", "missing", RemoveRequest{Hashcodes: []string{"h1"}}); !errors.Is(err, jobapi.ErrNotFound) {
		t.Errorf("Remove err = %v, want jobapi.ErrNotFound", err)
	}
}

func newTestJobString(t *testing.T, handler http.HandlerFunc) JobString {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := smclient.NewClient(server.Client(), "user", "secret")
	client.BaseURL = server.URL
	client.Credentials.AccessToken = &smclient.Token{
		Value:          "fake-test-token",
		ExpirationTime: time.Now().Add(1 * time.Hour),
	}

	return NewJobString(client)
}

func TestAdd_PostsHashcodesToAddPath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/strings/add") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got AddRequest
		_ = json.Unmarshal(body, &got)
		if len(got.Hashcodes) != 2 || got.Hashcodes[0] != "h1" {
			t.Errorf("hashcodes = %v, want [h1 h2]", got.Hashcodes)
		}
		if !got.MoveEnabled {
			t.Error("moveEnabled = false, want true")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"successCount":2,"failCount":0}}}`))
	}

	js := newTestJobString(t, handler)

	res, err := js.Add(context.Background(), "p1", "job-1", AddRequest{
		Hashcodes: []string{"h1", "h2"}, MoveEnabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SuccessCount != 2 || res.FailCount != 0 {
		t.Errorf("result = %+v, want {2 0}", res)
	}
}

func TestRemove_PostsHashcodesToRemovePath(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/strings/remove") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got RemoveRequest
		_ = json.Unmarshal(body, &got)
		if len(got.Hashcodes) != 1 || got.Hashcodes[0] != "h1" {
			t.Errorf("hashcodes = %v, want [h1]", got.Hashcodes)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{"successCount":1,"failCount":0}}}`))
	}

	js := newTestJobString(t, handler)

	res, err := js.Remove(context.Background(), "p1", "job-1", RemoveRequest{Hashcodes: []string{"h1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SuccessCount != 1 {
		t.Errorf("successCount = %d, want 1", res.SuccessCount)
	}
}

func TestList_ForwardsQueryAndMapsItems(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/jobs/job-1/strings") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("targetLocaleId") != "fr-FR" {
			t.Errorf("targetLocaleId = %q, want fr-FR", q.Get("targetLocaleId"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit = %q, want 10", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"code":"SUCCESS","data":{
			"totalCount": 2,
			"items": [
				{"targetLocaleId":"fr-FR","hashcode":"h1"},
				{"targetLocaleId":"fr-FR","hashcode":"h2"}
			]
		}}}`))
	}

	js := newTestJobString(t, handler)

	resp, err := js.List(context.Background(), "p1", "job-1", ListParams{TargetLocaleID: "fr-FR", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", resp.TotalCount)
	}
	if len(resp.Items) != 2 || resp.Items[0].Hashcode != "h1" || resp.Items[1].TargetLocaleID != "fr-FR" {
		t.Fatalf("unexpected items: %+v", resp.Items)
	}
}

func TestAddRemove_EmptyParamsRejectedBeforeRequest(t *testing.T) {
	called := false
	js := newTestJobString(t, func(http.ResponseWriter, *http.Request) { called = true })

	if _, err := js.Add(context.Background(), "p1", "job-1", AddRequest{}); !smerror.IsErrEmptyParam(err) {
		t.Errorf("Add with no hashcodes err = %v, want empty-param error", err)
	}
	if _, err := js.Remove(context.Background(), "", "job-1", RemoveRequest{Hashcodes: []string{"h1"}}); !smerror.IsErrEmptyParam(err) {
		t.Errorf("Remove with empty project err = %v, want empty-param error", err)
	}
	if _, err := js.List(context.Background(), "p1", "", ListParams{}); !smerror.IsErrEmptyParam(err) {
		t.Errorf("List with empty job err = %v, want empty-param error", err)
	}

	if called {
		t.Error("server was called despite empty-param validation")
	}
}
