package ai_platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	"go.temporal.io/sdk/client"
)

func (h *ZeusAiPlatformServicesWorker) ExecuteAiTaskWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []hermes_email_notifications.EmailContents) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiEmailWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, msgs)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTaskWorkflow")
		return err
	}
	return nil
}

type TelegramMessage struct {
	Timestamp   int    `json:"timestamp"`
	GroupName   string `json:"group_name"`
	SenderID    int    `json:"sender_id"`
	MessageText string `json:"message_text"`
	ChatID      int    `json:"chat_id"`
	MessageID   int    `json:"message_id"`
	TelegramMetadata
	//IsReply       bool   `json:"is_reply,omitempty"`
	//IsChannel     bool   `json:"is_channel,omitempty"`
	//IsGroup       bool   `json:"is_group,omitempty"`
	//IsPrivate     bool   `json:"is_private,omitempty"`
	//FirstName     string `json:"first_name,omitempty"`
	//LastName      string `json:"last_name,omitempty"`
	//Phone         string `json:"phone,omitempty"`
	//MutualContact bool   `json:"mutual_contact,omitempty"`
	//Username      string `json:"username,omitempty"`
}

type TelegramMetadata struct {
	IsReply       bool   `json:"is_reply,omitempty"`
	IsChannel     bool   `json:"is_channel,omitempty"`
	IsGroup       bool   `json:"is_group,omitempty"`
	IsPrivate     bool   `json:"is_private,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	MutualContact bool   `json:"mutual_contact,omitempty"`
	Username      string `json:"username,omitempty"`
}

func (h *ZeusAiPlatformServicesWorker) ExecuteAiTelegramWorkflow(ctx context.Context, ou org_users.OrgUser, msgs []TelegramMessage) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiIngestTelegramWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, msgs)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTaskWorkflow")
		return err
	}
	return nil
}
