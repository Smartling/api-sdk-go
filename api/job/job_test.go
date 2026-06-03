package job

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

func TestGetJob_PopulatesDetailFields(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"code": "SUCCESS",
				"data": {
					"translationJobUid": "jobUid-123",
					"jobName": "My Job",
					"jobNumber": "SMTL-7",
					"description": "desc",
					"referenceNumber": "ref-1",
					"jobStatus": "IN_PROGRESS",
					"dueDate": "2026-01-02T00:00:00Z",
					"createdDate": "2025-12-01T00:00:00Z",
					"modifiedDate": "2025-12-02T00:00:00Z",
					"createdByUserUid": "u-created",
					"modifiedByUserUid": "u-modified",
					"priority": 4,
					"targetLocaleIds": ["fr-FR"],
					"sourceFiles": [{"name": "a.json", "uri": "a.json", "fileUid": "f1"}],
					"customFields": [{"fieldUid": "cf1", "fieldName": "Dept", "fieldValue": "Eng"}],
					"issues": {"sourceIssuesCount": 1, "translationIssuesCount": 2}
				}
			}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.GetJob(context.Background(), "projectId-xyz", "jobUid-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.JobNumber != "SMTL-7" {
		t.Errorf("JobNumber = %q, want %q", resp.JobNumber, "SMTL-7")
	}
	if resp.JobStatus != "IN_PROGRESS" {
		t.Errorf("JobStatus = %q, want %q", resp.JobStatus, "IN_PROGRESS")
	}
	if resp.Priority != 4 {
		t.Errorf("Priority = %d, want 4", resp.Priority)
	}
	if len(resp.SourceFiles) != 1 || resp.SourceFiles[0].URI != "a.json" {
		t.Errorf("SourceFiles = %v, want one file a.json", resp.SourceFiles)
	}
	if len(resp.CustomFields) != 1 || resp.CustomFields[0].FieldName != "Dept" {
		t.Errorf("CustomFields = %v, want one Dept", resp.CustomFields)
	}
	if resp.Issues.TranslationIssuesCount != 2 {
		t.Errorf("Issues.TranslationIssuesCount = %d, want 2", resp.Issues.TranslationIssuesCount)
	}
}

func TestListProjectJobs_ForwardsFiltersAndMapsItems(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/projects/projectId-xyz/jobs") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("jobName") != "Release" {
			t.Errorf("jobName = %q, want Release", q.Get("jobName"))
		}
		if q.Get("translationJobStatus") != "IN_PROGRESS" {
			t.Errorf("translationJobStatus = %q, want IN_PROGRESS", q.Get("translationJobStatus"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit = %q, want 10", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {"code": "SUCCESS", "data": {
				"totalCount": 1,
				"items": [{"translationJobUid": "u1", "jobName": "Release", "jobStatus": "IN_PROGRESS", "targetLocaleIds": ["fr-FR"]}]
			}}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.ListProjectJobs(context.Background(), "projectId-xyz", ListProjectJobsParams{
		JobName:   "Release",
		JobStatus: []string{"IN_PROGRESS"},
		Page:      Page{Limit: 10},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", resp.TotalCount)
	}
	if len(resp.Items) != 1 || resp.Items[0].TranslationJobUID != "u1" {
		t.Fatalf("Items = %v, want one u1", resp.Items)
	}
	if resp.Items[0].JobStatus != "IN_PROGRESS" {
		t.Errorf("Items[0].JobStatus = %q, want IN_PROGRESS", resp.Items[0].JobStatus)
	}
}

func TestListAccountJobs_UsesAccountPathAndMapsProjectAndPriority(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/accounts/acct-1/jobs") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("withPriority") != "true" {
			t.Errorf("withPriority = %q, want true", r.URL.Query().Get("withPriority"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {"code": "SUCCESS", "data": {
				"totalCount": 1,
				"items": [{"translationJobUid": "u1", "jobName": "A", "projectId": "p9", "priority": 3}]
			}}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.ListAccountJobs(context.Background(), "acct-1", ListAccountJobsParams{WithPriority: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("Items = %v, want one", resp.Items)
	}
	if resp.Items[0].ProjectID != "p9" {
		t.Errorf("ProjectID = %q, want p9", resp.Items[0].ProjectID)
	}
	if resp.Items[0].Priority != 3 {
		t.Errorf("Priority = %d, want 3", resp.Items[0].Priority)
	}
}

func TestSearchJobs_PostsBodyAndMapsItems(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/p1/jobs/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string][]string
		_ = json.Unmarshal(body, &got)
		if len(got["fileUris"]) != 1 || got["fileUris"][0] != "a.json" {
			t.Errorf("fileUris = %v, want [a.json]", got["fileUris"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {"code": "SUCCESS", "data": {
				"totalCount": 1,
				"items": [{"translationJobUid": "u1", "jobName": "Found"}]
			}}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.SearchJobs(context.Background(), "p1", SearchJobsRequest{FileURIs: []string{"a.json"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].JobName != "Found" {
		t.Fatalf("Items = %v, want one Found", resp.Items)
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
