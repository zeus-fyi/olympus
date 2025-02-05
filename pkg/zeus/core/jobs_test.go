package zeus_core

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

func (s *JobsTestSuite) TestCreateRJob() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	var kns = zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "zeus",
		Env:           "production",
	}

	j := RedditJob("mlops")
	jc, err := s.K.CreateJob(ctx, kns, &j)
	s.Nil(err)
	s.Require().NotEmpty(jc)

}
func (s *JobsTestSuite) TestCreateJob() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
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

	// qn id
	/*
	   command: [ "/bin/sh" ]
	   args: [ "-ac"," /opt/web3signer/bin/web3signer --http-host-allowlist=* --http-listen-host=0.0.0.0 --tls-allow-any-client=true eth2 --key-manager-api-enabled=true --slashing-protection-db-url=${WEB3SIGNER_SLASHING_PROTECTION_DB_URL} --slashing-protection-db-password=${WEB3SIGNER_SLASHING_PROTECTION_DB_PASSWORD} --slashing-protection-db-username=${WEB3SIGNER_SLASHING_PROTECTION_DB_USERNAME} --network=${WEB3SIGNER_NETWORK}"]

	*/
	bof := int32(3)
	chID := "816017958395576400"

	timeAfter := time.Unix(int64(1700388338), 0).Add(-time.Hour * 3000).Format(time.RFC3339)

	fmt.Println(timeAfter)
	j := v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("discord-job-canary-%s", chID),
		},
		Spec: v1.JobSpec{
			BackoffLimit: &bof, // Setting backoffLimit to 0 to prevent retries
			Template: v1core.PodTemplateSpec{
				Spec: v1core.PodSpec{
					RestartPolicy: "OnFailure",
					InitContainers: []v1core.Container{
						{
							Name:    "discord-exporter-init",
							Image:   "tyrrrz/discordchatexporter:stable",
							Command: []string{"/bin/sh", "-ac"},
							Args: []string{
								fmt.Sprintf("/opt/app/DiscordChatExporter.Cli export -t %s --after \"%s\" -f Json -c %s -o /data/%s.json", authToken, timeAfter, chID, chID),
							},
							VolumeMounts: []v1core.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Containers: []v1core.Container{
						{
							Name:            "discord-job",
							Image:           "zeusfyi/snapshots:latest",
							ImagePullPolicy: "Always",
							Command:         []string{"/bin/sh", "-c"},
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
	err = s.K.DeleteJob(ctx, kns, j.Name)
	s.Nil(err)
	err = s.K.DeleteAllPodsLike(ctx, kns, j.Name, nil, nil)
	s.Nil(err)

	// "discord-exporter-init "
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

func RedditJob(subreddit string) v1.Job {
	bof := int32(2)
	j := v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("reddit-job-%s", subreddit),
		},
		Spec: v1.JobSpec{
			BackoffLimit: &bof, // Setting backoffLimit to 0 to prevent retries
			Template: v1core.PodTemplateSpec{
				Spec: v1core.PodSpec{
					//ImagePullSecrets: []v1core.LocalObjectReference{
					//	{
					//		Name: "zeus-fyi-ext",
					//	},
					//},
					RestartPolicy: "OnFailure",
					Containers: []v1core.Container{
						{
							Name:            "reddit-job",
							Image:           "registry.digitalocean.com/zeus-fyi/hephaestus:latest",
							ImagePullPolicy: "Always",
							Command:         []string{"/bin/sh", "-c"},
							Args: []string{
								fmt.Sprintf("exec hephaestus --workload-type=\"%s\"", subreddit),
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
	return j
}
