package zeus_core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobsTestSuite struct {
	K8TestSuite
}

var ctx = context.Background()

func (s *JobsTestSuite) TestGetJobsList() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ephemeral"
	jl, err := s.K.GetJobsList(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(jl)
}

const (
	internalUser = 7138958574876245567
)

func (s *JobsTestSuite) TestCreateJob() {
	var kns = zeus_common_types.CloudCtxNs{
		CloudProvider: "ovh",
		Region:        "us-west-or-1",
		Context:       "kubernetes-admin@zeusfyi",
		Namespace:     "zeus",
		Env:           "production",
	}
	authToken, err := read_keys.GetDiscordKey(ctx, internalUser)
	s.Require().Nil(err)

	hs, err := misc.HashParams([]interface{}{authToken})
	s.Require().Nil(err)

	chID := "7138958574876245567"
	ts := time.Now().Add(-time.Hour).Format(time.RFC3339)
	j := v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "discord-job",
		},
		Spec: v1.JobSpec{
			Template: v1core.PodTemplateSpec{
				Spec: v1core.PodSpec{
					InitContainers: []v1core.Container{
						{
							Name:    "discord-exporter-init",
							Image:   "tyrrrz/discordchatexporter:stable",
							Command: []string{"sh", "-c"},
							Args: []string{
								fmt.Sprintf("discordchatexporter export -t %s --after \"%s\" -f Json -c %s -o /data/%s.json", authToken, ts, chID, chID),
							},
						},
					},
					Containers: []v1core.Container{
						{
							Name:    "discord-job",
							Image:   "zeusfyi/snapshots:latest",
							Command: []string{"sh", "-c"},
							Args: []string{
								fmt.Sprintf("exec snapshots --bearer=\"%s\" --payload-base-path=\"https://api.zeus.fyi\" --payload-post-path=\"/vz/webhooks/discord/ai\" --workload-type=\"send-payload\" --fi %s.json", hs, chID),
							},
							VolumeMounts: []v1core.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []v1core.Volume{
						{
							Name: "data-volume",
							VolumeSource: v1core.VolumeSource{
								EmptyDir: &v1core.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	jc, err := s.K.CreateJob(ctx, kns, &j)
	s.Nil(err)
	s.Require().NotEmpty(jc)
}

func (s *JobsTestSuite) TestDeleteJob() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	err := s.K.DeleteJob(ctx, kns, "jobName")
	s.Nil(err)
}

func TestJobsTestSuite(t *testing.T) {
	suite.Run(t, new(JobsTestSuite))
}
