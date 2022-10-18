package create

import (
	common_fields "github.com/zeus-fyi/olympus/pkg/hera/cookbook/recipes/common/fields"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/fields"
)

func genInsertFnFields() []fields.Field {
	return []fields.Field{common_fields.CtxField(), common_fields.QueryParams()}
}

func genInsertFnReturnFields() []fields.Field {
	return []fields.Field{common_fields.ErrField()}
}
