package mt

import "github.com/Smartling/api-sdk-go/helpers/uid"

// UploadFileResponse defines upload file response
type UploadFileResponse struct {
	Code    int
	FileUID uid.FileUID
}

// UploadFileRequest defines upload file request
type UploadFileRequest struct {
	File               []byte
	FileType           Type
	LocalesToAuthorize []string
	Directives         map[string]string
}

type uploadFileResponse struct {
	Response struct {
		Code int
		Data struct {
			FileUID string `json:"fileUid"`
		} `json:"data"`
	} `json:"response"`
}

func toUploadFileResponse(r uploadFileResponse) UploadFileResponse {
	return UploadFileResponse{
		Code:    r.Response.Code,
		FileUID: uid.FileUID(r.Response.Data.FileUID),
	}
}
