// methods that help with encoding structures to GET url queries

package smartling

import (
	"net/url"
	"strconv"
	"log"
)

func (plr *ProjectListRequest) RawQuery() string {
	q := url.Values{}
	if len(plr.ProjectNameFilter) > 0 {
		q.Set("projectNameFilter", plr.ProjectNameFilter)
	}
	q.Set("includeArchived", strconv.FormatBool(plr.IncludeArchived))
	if plr.Limit > 0 {
		q.Set("limit", strconv.FormatInt(plr.Limit, 10))
	}
	q.Set("offset", strconv.FormatInt(plr.Offset, 10))

	return q.Encode()
}

func (flr *FileListRequest) RawQuery() string {
	q := url.Values{}
	if len(flr.UriMask) > 0 {
		q.Set("uriMask", flr.UriMask)
	}
	for _, tp := range flr.FileTypes {
		q.Add("fileTypes[]", string(tp))
	}
	if !flr.LastUploadedAfter.IsZero() {
		log.Printf("not a zero")
		flr.LastUploadedAfter.EncodeValues("lastUploadedAfter", &q)
	}
	if !flr.LastUploadedBefore.IsZero() {
		flr.LastUploadedBefore.EncodeValues("lastUploadedBefore", &q)
	}
	if flr.Limit > 0 {
		q.Set("limit", strconv.FormatInt(flr.Limit, 10))
	}
	q.Set("offset", strconv.FormatInt(flr.Offset, 10))

	return q.Encode()
}
