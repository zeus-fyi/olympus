package fields

import (
	"github.com/zeus-fyi/tables-to-go/pkg/database"
)

type DbMetadata struct {
	*database.Table
	database.Column
}

func NewDbMetadata(t *database.Table, c database.Column) DbMetadata {
	return DbMetadata{t, c}
}
