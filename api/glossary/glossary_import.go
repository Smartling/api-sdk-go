package glossary

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

const (
	PendingImportStatus    = "PENDING"
	InProgressImportStatus = "IN_PROGRESS"
	SuccessfulImportStatus = "SUCCESSFUL"
	FailedImportStatus     = "FAILED"
)

// Import uploads a glossary file to the given glossary as a multipart request.
// Endpoint: POST /glossary-api/v3/accounts/{accountUid}/glossaries/{glossaryUid}/import.
func (h httpGlossary) Import(ctx context.Context,
	accountUID uid.AccountUID, glossaryUID string, req ImportGlossaryRequest) (ImportGlossaryResponse, error) {
	body, contentType, err := buildImportForm(req)
	if err != nil {
		return ImportGlossaryResponse{}, fmt.Errorf("failed to build import form: %w", err)
	}

	reqURL := path.Join(glossaryBasePath, url.PathEscape(string(accountUID)), "glossaries", url.PathEscape(glossaryUID), "import")

	var response importGlossary
	_, code, err := h.client.PostJSON(
		ctx,
		reqURL,
		body,
		&response.Response.Data,
		smclient.ContentTypeOption(contentType),
	)
	if err != nil && code == http.StatusNotFound {
		return ImportGlossaryResponse{}, ErrGlossaryNotFound
	}
	if err != nil {
		return ImportGlossaryResponse{}, fmt.Errorf("failed to import glossary: %w", err)
	}

	return toImportGlossaryResponse(response, code), nil
}

// ImportStatus polls the status of a previously submitted glossary import.
// Endpoint: GET /glossary-api/v3/accounts/{accountUid}/glossaries/{glossaryUid}/import/{importUid}.
func (h httpGlossary) ImportStatus(ctx context.Context, accountUID uid.AccountUID, glossaryUID, importUID string) (ImportStatusResponse, error) {
	reqURL := path.Join(
		glossaryBasePath,
		url.PathEscape(string(accountUID)),
		"glossaries", url.PathEscape(glossaryUID),
		"import", url.PathEscape(importUID),
	)

	var response importStatusResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &response.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return ImportStatusResponse{}, ErrImportNotFound
	}
	if err != nil {
		return ImportStatusResponse{}, fmt.Errorf("failed to get import status: %w", err)
	}

	return toImportStatusResponse(response, code), nil
}

// ImportConfirm confirms a previously created glossary import.
// Only imports in PENDING status may be confirmed.
// Endpoint: POST /glossary-api/v3/accounts/{accountUid}/glossaries/{glossaryUid}/import/{importUid}/confirm.
func (h httpGlossary) ImportConfirm(ctx context.Context, accountUID uid.AccountUID, glossaryUID, importUID string) (bool, error) {
	reqURL := path.Join(glossaryBasePath, url.PathEscape(string(accountUID)), "glossaries", url.PathEscape(glossaryUID), "import", url.PathEscape(importUID), "confirm")

	_, code, err := h.client.PostJSON(ctx, reqURL, nil, nil)
	if err != nil && code == http.StatusNotFound {
		return false, ErrImportNotFound
	}
	if err != nil {
		return false, fmt.Errorf("failed to confirm glossary import: %w", err)
	}

	return true, nil
}

func buildImportForm(req ImportGlossaryRequest) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Per Smartling/api-docs GlossaryImportCommandPTO, the form fields are:
	// importFile (binary), importFileName, importFileMediaType, archiveMode.
	filePart, err := writer.CreateFormFile("importFile", req.FileName)
	if err != nil {
		return nil, "", fmt.Errorf("create importFile part: %w", err)
	}
	if _, err := filePart.Write(req.File); err != nil {
		return nil, "", fmt.Errorf("write importFile part: %w", err)
	}

	if err := writer.WriteField("importFileName", req.FileName); err != nil {
		return nil, "", fmt.Errorf("write importFileName field: %w", err)
	}
	if err := writer.WriteField("importFileMediaType", req.MediaType); err != nil {
		return nil, "", fmt.Errorf("write importFileMediaType field: %w", err)
	}
	if err := writer.WriteField("archiveMode", strconv.FormatBool(req.ArchiveMode)); err != nil {
		return nil, "", fmt.Errorf("write archiveMode field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart writer: %w", err)
	}

	return body.Bytes(), writer.FormDataContentType(), nil
}
