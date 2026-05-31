package glossary

type ImportStatusResponse struct {
	Code         int
	GlossaryUID  string
	ImportUID    string
	ImportStatus string
}

func toImportStatusResponse(r importStatusResponse, code int) ImportStatusResponse {
	return ImportStatusResponse{
		Code:         code,
		GlossaryUID:  r.Response.Data.GlossaryUID,
		ImportUID:    r.Response.Data.ImportUID,
		ImportStatus: r.Response.Data.ImportStatus,
	}
}

type importStatusResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			GlossaryUID  string `json:"glossaryUid"`
			ImportUID    string `json:"importUid"`
			ImportStatus string `json:"importStatus"`
		} `json:"data"`
	} `json:"response"`
}
