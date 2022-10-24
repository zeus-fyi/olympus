package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/common"
)

func (d *Deployment) insertSpecChildren(parentSqlExpression, cteSpecParent string) string {
	s := d.Spec
	// should be three child types replica, selector, template (pod/spec)

	// replica TODO verify
	parentSqlExpression = common.InsertChildClassSingleValueType(parentSqlExpression, cteSpecParent, s.Replicas)

	// selector
	parentSqlExpression = common.InsertChildClassValues(parentSqlExpression, cteSpecParent, s.Selector.MatchLabels)

	// template
	parentSqlExpression = d.insertPodSpecTemplate(parentSqlExpression, cteSpecParent)
	return parentSqlExpression
}
