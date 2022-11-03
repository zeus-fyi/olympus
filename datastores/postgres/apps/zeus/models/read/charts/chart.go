package read_charts

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	read_configuration "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/configuration"
	read_deployments "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/deployments"
	read_networking "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/networking"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Chart struct {
	charts.Chart

	chart_workload.ChartWorkload
}

func NewChartReader() Chart {
	c := charts.NewChart()
	k8s := chart_workload.NewChartWorkload()
	cr := Chart{
		Chart:         c,
		ChartWorkload: k8s,
	}
	return cr
}

const ModelName = "Chart"

func (c *Chart) SelectSingleChartsResources(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("SelectQuery", q.LogHeader(ModelName))

	rows, err := apps.Pg.Query(ctx, q.RawQuery, q.CTEQuery.Params...)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(ModelName))
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var podSpecVolumesStr string
		container := containers.NewContainer()

		var ckagg string
		rowErr := rows.Scan(&c.ChartPackageID, &c.ChartName, &c.ChartVersion, &c.ChartDescription, &c.ChartComponentKindName,
			&ckagg,
			&container.Metadata.ContainerID, &container.Metadata.ContainerName, &container.Metadata.ContainerImageID,
			&container.Metadata.ContainerVersionTag, &container.Metadata.ContainerPlatformOs, &container.Metadata.ContainerRepository,
			&container.Metadata.ContainerImagePullPolicy,
			&container.DB.Ports,
			&container.DB.EnvVar, &container.DB.Probes, &container.DB.ContainerVolumes, &container.DB.CmdArgs,
			&podSpecVolumesStr,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return rowErr
		}
		switch c.ChartComponentKindName {
		case "Deployment":
			if c.Deployment == nil {
				deployment := deployments.NewDeployment()
				c.Deployment = &deployment
				derr := read_deployments.DBDeploymentResource(c.Deployment, ckagg, podSpecVolumesStr)
				if derr != nil {
					log.Err(derr).Msg(q.LogHeader(ModelName))
					return derr
				}
			}
			err = read_deployments.DBDeploymentContainer(c.Deployment, &container)
			if err != nil {
				log.Err(err).Msg(q.LogHeader(ModelName))
				return err
			}
		case "Service":
			if c.Service == nil {
				svc := services.NewService()
				c.Service = &svc
				serr := read_networking.DBServiceResource(c.Service, ckagg)
				if serr != nil {
					log.Err(serr).Msg(q.LogHeader(ModelName))
					return serr
				}
			}
		case "Ingress":
			if c.Ingress == nil {
				ing := ingresses.NewIngress()
				c.Ingress = &ing
				ierr := read_networking.DBIngressResource(c.Ingress, ckagg)
				if ierr != nil {
					log.Err(ierr).Msg(q.LogHeader(ModelName))
					return ierr
				}
			}
		case "ConfigMap":
			if c.ConfigMap == nil {
				cm := configuration.NewConfigMap()
				c.ConfigMap = &cm
				cerr := read_configuration.DBConfigMapResource(c.ConfigMap, ckagg)
				if cerr != nil {
					log.Err(cerr).Msg(q.LogHeader(ModelName))
					return cerr
				}
			}
		}
	}
	return nil
}
