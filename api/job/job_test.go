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
	wantDue := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	if !resp.Dates.Due.Equal(wantDue) {
		t.Errorf("Dates.Due = %v, want %v", resp.Dates.Due, wantDue)
	}
	wantCreated := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	if !resp.Dates.Created.Equal(wantCreated) {
		t.Errorf("Dates.Created = %v, want %v", resp.Dates.Created, wantCreated)
	}
	wantModified := time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC)
	if !resp.Dates.Modified.Equal(wantModified) {
		t.Errorf("Dates.Modified = %v, want %v", resp.Dates.Modified, wantModified)
	}
	if !resp.Dates.FirstCompleted.IsZero() {
		t.Errorf("Dates.FirstCompleted = %v, want zero (job in progress)", resp.Dates.FirstCompleted)
	}
	if !resp.Dates.LastCompleted.IsZero() {
		t.Errorf("Dates.LastCompleted = %v, want zero (job in progress)", resp.Dates.LastCompleted)
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

func TestFindJobsByStrings_PostsBodyAndMapsItems(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/projects/p1/jobs/find-jobs-by-strings") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string][]string
		_ = json.Unmarshal(body, &got)
		if !equalStringSlices(got["hashcodes"], []string{"h1"}) {
			t.Errorf("hashcodes = %v, want [h1]", got["hashcodes"])
		}
		if !equalStringSlices(got["localeIds"], []string{"fr-FR"}) {
			t.Errorf("localeIds = %v, want [fr-FR]", got["localeIds"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {"code": "SUCCESS", "data": {
				"totalCount": 1,
				"items": [{
					"translationJobUid": "u1",
					"jobName": "Found",
					"dueDate": null,
					"hashcodesByLocale": [
						{"localeId": "fr-FR", "hashcodes": ["h1"]},
						{"localeId": "de-DE", "hashcodes": ["h1"]}
					]
				}]
			}}
		}`))
	}

	j, _ := newTestJob(t, handler)

	resp, err := j.FindJobsByStrings(context.Background(), "p1", FindJobsByStringsRequest{
		Hashcodes: []string{"h1"},
		LocaleIDs: []string{"fr-FR"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 || len(resp.Items) != 1 {
		t.Fatalf("resp = %+v, want totalCount 1 / one item", resp)
	}
	item := resp.Items[0]
	if item.TranslationJobUID != "u1" || item.JobName != "Found" {
		t.Errorf("item = %+v, want u1/Found", item)
	}
	if len(item.HashcodesByLocale) != 2 {
		t.Fatalf("HashcodesByLocale = %+v, want 2 locales", item.HashcodesByLocale)
	}
	if item.HashcodesByLocale[0].LocaleID != "fr-FR" ||
		!equalStringSlices(item.HashcodesByLocale[0].Hashcodes, []string{"h1"}) {
		t.Errorf("first locale = %+v, want fr-FR/[h1]", item.HashcodesByLocale[0])
	}
}

func TestFindJobsByStrings_DueDateParsing(t *testing.T) {
	tests := []struct {
		name    string
		dueDate string
		wantErr bool
		want    time.Time
	}{
		{
			name:    "valid RFC3339 date is parsed",
			dueDate: `"2015-11-21T11:51:17Z"`,
			want:    time.Date(2015, 11, 21, 11, 51, 17, 0, time.UTC),
		},
		{
			name:    "malformed date returns error",
			dueDate: `"not-a-date"`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"response": {"code": "SUCCESS", "data": {
						"totalCount": 1,
						"items": [{
							"translationJobUid": "u1",
							"jobName": "Found",
							"dueDate": ` + tt.dueDate + `,
							"hashcodesByLocale": [{"localeId": "fr-FR", "hashcodes": ["h1"]}]
						}]
					}}
				}`))
			}

			j, _ := newTestJob(t, handler)

			resp, err := j.FindJobsByStrings(context.Background(), "p1", FindJobsByStringsRequest{
				Hashcodes: []string{"h1"},
			})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !resp.Items[0].DueDate.Equal(tt.want) {
				t.Errorf("DueDate = %v, want %v", resp.Items[0].DueDate, tt.want)
			}
		})
	}
}

func TestFindJobsByStrings_NotFoundMapsToErrNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	j, _ := newTestJob(t, handler)

	_, err := j.FindJobsByStrings(context.Background(), "p1", FindJobsByStringsRequest{
		Hashcodes: []string{"h1"},
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
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
