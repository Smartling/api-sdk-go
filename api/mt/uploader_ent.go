package mt

// UploadFileResponse defines upload file response
type UploadFileResponse struct {
	Code    string
	FileUID FileUID
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
		Code string `json:"code"`
		Data struct {
			FileUID string `json:"fileUid"`
		} `json:"data"`
	} `json:"response"`
}

func toUploadFileResponse(r uploadFileResponse) UploadFileResponse {
	return UploadFileResponse{
		Code:    r.Response.Code,
		FileUID: FileUID(r.Response.Data.FileUID),
	}
}
