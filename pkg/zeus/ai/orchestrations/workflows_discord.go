package ai_platform_service_orchestrations

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (z *ZeusAiPlatformServiceWorkflows) AiIngestDiscordWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, searchGroupName string, cm hera_discord.ChannelMessages) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	if len(cm.Messages) == 0 {
		return nil
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiIngestDiscordWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	var sq *hera_search.DiscordSearchResultWrapper

	searchQueryCtx := workflow.WithActivityOptions(ctx, ao)
	if len(cm.Guild.Id) > 0 && len(cm.Channel.Id) > 0 {
		err = workflow.ExecuteActivity(searchQueryCtx, z.SelectDiscordSearchQueryByGuildChannel, ou, cm.Guild.Id, cm.Channel.Id).Get(searchQueryCtx, &sq)
		if err != nil {
			logger.Error("failed to execute SelectDiscordSearchQuery", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
	} else {
		err = workflow.ExecuteActivity(searchQueryCtx, z.SelectDiscordSearchQuery, ou, searchGroupName).Get(searchQueryCtx, &sq)
		if err != nil {
			logger.Error("failed to execute SelectDiscordSearchQuery", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
	}

	if sq == nil && sq.SearchID == 0 {
		logger.Info("no search id found")
		return nil
	}
	insertMessagesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(insertMessagesCtx, z.InsertIncomingDiscordDataFromSearch, sq.SearchID, cm).Get(insertMessagesCtx, nil)
	if err != nil {
		logger.Error("failed to execute InsertIncomingDiscordDataFromSearch", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

func (z *ZeusAiPlatformServiceWorkflows) AiFetchDataToIngestDiscordWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, searchGroupName string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiFetchDataToIngestDiscordWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	searchQueryCtx := workflow.WithActivityOptions(ctx, ao)
	var sq *hera_search.DiscordSearchResultWrapper
	err = workflow.ExecuteActivity(searchQueryCtx, z.SelectDiscordSearchQuery, ou, searchGroupName).Get(searchQueryCtx, &sq)
	if err != nil {
		logger.Error("failed to execute SelectDiscordSearchQuery", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	if sq == nil || sq.SearchID == 0 {
		logger.Info("no new tweets found")
		return nil
	}
	for i, jib := range sq.Results {
		if i > 0 {
			err = workflow.Sleep(ctx, time.Second*5)
			if err != nil {
				logger.Error("failed to sleep", "Error", err)
				return err
			}
		}
		jobCtx := workflow.WithActivityOptions(ctx, ao)
		timeAfter := time.Unix(int64(jib.MaxMessageID), 0).Add(-time.Minute * 5).Format(time.RFC3339)
		err = workflow.ExecuteActivity(jobCtx, z.CreateDiscordJob, ou, sq.SearchID, jib.ChannelID, timeAfter).Get(jobCtx, nil)
		if err != nil {
			logger.Error("failed to execute CreateDiscordJob", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

func DiscordJob(orgID, si int, authToken, hs, chID, ts string) v1.Job {
	bof := int32(0)
	j := v1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("d-job-%d-%s", si, chID),
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
								fmt.Sprintf("/opt/app/DiscordChatExporter.Cli export -t %s --after \"%s\" -f Json -c %s -o /data/%s.json", authToken, ts, chID, chID),
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
								fmt.Sprintf("exec snapshots --bearer=\"%s\" --payload-base-path=\"https://api.zeus.fyi\" --payload-post-path=\"/vz/webhooks/discord/ai/%d\" --workload-type=\"send-payload\" --fi %s.json", hs, orgID, chID),
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

func RedditJob(subreddit string) v1.Job {
	bof := int32(0)
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
					RestartPolicy: "OnFailure",
					Containers: []v1core.Container{
						{
							Name:            "reddit-job",
							Image:           "zeusfyi/hephaestus:latest",
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
