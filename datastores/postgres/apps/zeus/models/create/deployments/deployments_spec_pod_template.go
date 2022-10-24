package deployments

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/common"
)

// TODO refactor to get rid of parentSqlExpression to use other style
func (d *Deployment) insertPodSpecTemplate(parentSqlExpression, parentClassTypeCteName string) string {
	pts := d.Spec.Template.Spec

	// this part links containers to this k8s workload resource
	parentSqlExpression = d.Spec.DeploymentSpec.Template.InsertPodContainerGroupSQL() + parentSqlExpression
	// TODO, maybe not even rely on this, and have the containers already exist...
	// add containers now, this func will also get the other data from the containers loop iterator
	parentSqlExpression = common.InsertContainerValues(parentSqlExpression, pts.PodTemplateContainers)

	// optionally add template metadata here if desired later
	return parentSqlExpression
}
