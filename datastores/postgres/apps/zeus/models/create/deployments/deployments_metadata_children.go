package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/common"
)

func (d *Deployment) insertDeploymentMetadataChildren(parentSqlExpression, cteMetadataParent string) string {
	m := d.Metadata
	if m.HasName() {
		parentSqlExpression = common.InsertChildClassSingleValueType(parentSqlExpression, cteMetadataParent, m.Name)
	}
	if m.HasLabels() {
		parentSqlExpression = common.InsertChildClassValues(parentSqlExpression, cteMetadataParent, m.Labels)
	}
	if m.HasAnnotations() {
		parentSqlExpression = common.InsertChildClassValues(parentSqlExpression, cteMetadataParent, m.Annotations)
	}
	return parentSqlExpression
}
