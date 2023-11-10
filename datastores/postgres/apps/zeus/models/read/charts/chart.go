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
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/jobs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/servicemonitors"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
	read_configuration "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/configuration"
	read_deployments "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/deployments"
	read_jobs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/jobs"
	read_networking "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/networking"
	read_servicemonitors "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/servicemonitors"
	read_statefulsets "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/statefulsets"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Chart struct {
	charts.Chart `json:"chart"`

	chart_workload.ChartWorkload `json:"chartWorkload"`
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
			&container.Metadata.ContainerImagePullPolicy, &container.Metadata.IsInitContainer,
			&container.DB.Ports, &container.DB.EnvVar, &container.DB.Probes, &container.DB.ContainerVolumes,
			&container.DB.CmdArgs, &container.DB.ComputeResources, &container.DB.SecurityContext,
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
		case "StatefulSet":
			if c.StatefulSet == nil {
				sts := statefulsets.NewStatefulSet()
				c.StatefulSet = &sts
				derr := read_statefulsets.DBStatefulSetResource(c.StatefulSet, ckagg, podSpecVolumesStr)
				if derr != nil {
					log.Err(derr).Msg(q.LogHeader(ModelName))
					return derr
				}
			}
			err = read_statefulsets.DBStatefulSetContainer(c.StatefulSet, &container)
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
		case "ServiceMonitor":
			if c.ServiceMonitor == nil {
				sm := servicemonitors.NewServiceMonitor()
				c.ServiceMonitor = &sm
				serr := read_servicemonitors.DBServiceMonitorResource(c.ServiceMonitor, ckagg)
				if serr != nil {
					log.Err(serr).Msg(q.LogHeader(ModelName))
					return serr
				}
			}
		case "Job":
			if c.Job == nil {
				j := jobs.NewJob()
				c.Job = &j
				jerr := read_jobs.DBJobResource(c.Job, ckagg)
				if jerr != nil {
					log.Err(jerr).Msg(q.LogHeader(ModelName))
					return jerr
				}
			}
		case "CronJob":
			if c.CronJob == nil {
				cj := jobs.NewCronJob()
				c.CronJob = &cj
				cjerr := read_jobs.DBCronJobResource(c.CronJob, ckagg)
				if cjerr != nil {
					log.Err(cjerr).Msg(q.LogHeader(ModelName))
					return cjerr
				}
			}
		}
	}
	return nil
}
