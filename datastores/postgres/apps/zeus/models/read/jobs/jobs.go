package read_jobs

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/jobs"
)

func DBJobResource(j *jobs.Job, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		log.Err(pcerr)
		return pcerr
	}
	err := j.ParseDBConfigToK8s(pcGroupMap)
	return err
}

func DBCronJobResource(cj *jobs.CronJob, ckagg string) error {
	pcGroupMap, pcerr := common_conversions.ParseParentChildAggValues(ckagg)
	if pcerr != nil {
		log.Err(pcerr)
		return pcerr
	}
	err := cj.ParseDBConfigToK8s(pcGroupMap)
	return err
}
