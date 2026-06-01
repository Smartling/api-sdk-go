package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// Export downloads a glossary as a binary file (TBX/CSV/XLSX, etc.).
// The format is determined by req.Format. On success the response Data is the
// open HTTP response body - the caller MUST Close it.
func (h httpGlossary) Export(ctx context.Context,
	accountUID uid.AccountUID, glossaryUID string, req ExportGlossaryRequest) (ExportGlossaryResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return ExportGlossaryResponse{}, fmt.Errorf("failed to marshal export request: %w", err)
	}

	reqURL := path.Join(glossaryBasePath, string(accountUID), "glossaries", glossaryUID, "entries", "download")

	reply, err := h.client.Post(ctx, reqURL, payload)
	if err != nil {
		return ExportGlossaryResponse{}, fmt.Errorf("failed to export glossary: %w", err)
	}

	if reply.StatusCode == http.StatusNotFound {
		if err := reply.Body.Close(); err != nil {
			h.client.Logger.Infof("failed to close response body: %v", err)
		}
		return ExportGlossaryResponse{}, ErrGlossaryNotFound
	}
	if reply.StatusCode != http.StatusOK {
		body, err := io.ReadAll(reply.Body)
		if err != nil {
			h.client.Logger.Infof("failed to read response body: %v", err)
		}
		if err := reply.Body.Close(); err != nil {
			h.client.Logger.Infof("failed to close response body: %v", err)
		}
		return ExportGlossaryResponse{}, fmt.Errorf(
			"failed to export glossary: unexpected response code %d: %s",
			reply.StatusCode, body,
		)
	}

	contentType := reply.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		body, err := io.ReadAll(reply.Body)
		if err != nil {
			h.client.Logger.Infof("failed to read response body: %v", err)
		}
		if err := reply.Body.Close(); err != nil {
			h.client.Logger.Infof("failed to close response body: %v", err)
		}
		return ExportGlossaryResponse{}, fmt.Errorf(
			"failed to export glossary: server returned JSON instead of file: %s",
			body,
		)
	}

	h.client.Logger.Debugf(
		"export response headers: Content-Type=%q Content-Length=%d Content-Disposition=%q Location=%q",
		contentType,
		reply.ContentLength,
		reply.Header.Get("Content-Disposition"),
		reply.Header.Get("Location"),
	)

	return ExportGlossaryResponse{
		Code:          reply.StatusCode,
		Filename:      filenameFromContentDisposition(reply.Header.Get("Content-Disposition")),
		ContentType:   contentType,
		ContentLength: reply.ContentLength,
		Data:          reply.Body,
	}, nil
}

func filenameFromContentDisposition(header string) string {
	if header == "" {
		return ""
	}
	if _, params, err := mime.ParseMediaType(header); err == nil {
		return params["filename"]
	}
	return ""
}
