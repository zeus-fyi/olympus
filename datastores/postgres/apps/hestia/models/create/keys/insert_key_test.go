package create_keys

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateKeyTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateKeyTestSuite) TestInsertUserKey() {
	ctx := context.Background()

	uID := s.NewTestUser()
	nk := NewCreateKey(uID, "0x974C0c36265b7aa658b63A6121041AeE9e4DFd1b")
	nk.PublicKeyVerified = false
	nk.PublicKeyName = "test_key"
	nk.PublicKeyTypeID = keys.EcdsaKeyTypeID
	nk.CreatedAt = time.Now()
	q := sql_query_templates.NewQueryParam("InsertUserKey", "user_keys", "where", 1000, []string{})
	q.TableName = nk.GetTableName()
	q.Columns = nk.GetTableColumns()
	q.Values = []apps.RowValues{nk.GetRowValues("default")}

	err := nk.InsertUserKey(ctx, q)
	s.Require().Nil(err)
}

func TestCreateKeyTestSuite(t *testing.T) {
	suite.Run(t, new(CreateKeyTestSuite))
}
