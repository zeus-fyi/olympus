package deployments

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/common"

func (d *Deployment) insertPodSpecTemplate(parentSqlExpression, parentClassTypeCteName string) string {
	pts := d.Spec.Template.Spec

	// TODO, maybe not even rely on this, and have the containers already exist...
	// add containers now, this func will also get the other data from the containers loop iterator
	parentSqlExpression = common.InsertContainerValues(parentSqlExpression, pts.PodTemplateContainers)

	// optionally add template metadata here if desired later

	// this part links containers to this k8s workload resource
	// TODO disabling for now, need to get cont ids
	// parentSqlExpression = common.SetPodSpecTemplateChildTypeInsert(parentSqlExpression, parentClassTypeCteName, pts.PodTemplateSpecClassDefinition)

	return parentSqlExpression
}
