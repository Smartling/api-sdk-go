package glossary

import (
	"context"
	"errors"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

const glossaryBasePath = "/glossary-api/v3/accounts/"

var (
	ErrGlossaryNotFound = errors.New("glossary not found")
	ErrImportNotFound   = errors.New("glossary import not found")
)

// Glossary defines the glossary behaviour
type Glossary interface {
	Get(ctx context.Context, accountUID uid.AccountUID, glossaryUID string) (GetGlossaryResponse, error)
	GetByName(ctx context.Context, accountUID uid.AccountUID, name string) (glossaries []GetGlossaryResponse, err error)
	Import(ctx context.Context, accountUID uid.AccountUID, glossaryUID string, req ImportGlossaryRequest) (ImportGlossaryResponse, error)
	ImportStatus(ctx context.Context, accountUID uid.AccountUID, glossaryUID, importUID string) (ImportStatusResponse, error)
	ImportConfirm(ctx context.Context, accountUID uid.AccountUID, glossaryUID, importUID string) (bool, error)
	Export(ctx context.Context, accountUID uid.AccountUID, glossaryUID string, req ExportGlossaryRequest) (ExportGlossaryResponse, error)
	Create(ctx context.Context, accountUID uid.AccountUID, req CreateGlossaryRequest) (CreateGlossaryResponse, error)
}

// NewGlossary returns new Glossary implementation
func NewGlossary(client *smclient.Client) Glossary {
	return newHttpGlossary(client)
}

// httpGlossary implements Glossary interface using HTTP client
type httpGlossary struct {
	client *smclient.Client
}

func newHttpGlossary(client *smclient.Client) httpGlossary {
	return httpGlossary{client: client}
}
