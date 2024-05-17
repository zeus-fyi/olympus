package zeus_v1_ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type FlowsActionsGetRequest struct{}

func FlowsExportCsvRequestHandler(c echo.Context) error {
	request := new(FlowsActionsGetRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}
	return ExportRunCsvRequest(c, id)
}

func ExportRunCsvRequest(c echo.Context, id int) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ojsRuns, err := artemis_orchestrations.SelectAiSystemOrchestrations(c.Request().Context(), ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if len(ojsRuns) < 1 {
		return c.JSON(http.StatusOK, nil)
	}
	ojr := ojsRuns[0]
	// 7138983863666903883
	//ojr := artemis_orchestrations.OrchestrationsAnalysis{}
	//ojr.OrchestrationName = "test-wf"
	log.Info().Interface("oj.GroupName", ojr.GroupName).Msg("ExportRunCsvRequest")
	ue := artemis_entities.UserEntity{
		Nickname: ojr.OrchestrationName,
		Platform: "csv-exports",
	}
	_, err = ai_platform_service_orchestrations.S3WfRunExport(c.Request().Context(), ou, ojr.OrchestrationName, &ue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	for _, v := range ue.MdSlice {
		if v.TextData == nil {
			continue
		}
		c.Response().Header().Set("Content-Type", "text/csv")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, ojr.GroupName))
		_, err = c.Response().Write([]byte(*v.TextData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return nil
}

func AdminFlowsExportCsvRequestHandler(c echo.Context) error {
	request := new(FlowsActionsGetRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusBadRequest, "invalid ID parameter")
	}
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if ou.OrgID != 1710298581127603000 && ou.OrgID != 7138983863666903883 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	ueSlice, err := ExportRunCsvRequest2(c.Request().Context(), ou, id)
	if err != nil {
		log.Err(err).Msg("invalid ID parameter")
		return c.JSON(http.StatusInternalServerError, err)
	}
	for _, ue := range ueSlice {
		for _, v := range ue.MdSlice {
			if v.TextData == nil {
				continue
			}
			c.Response().Header().Set("Content-Type", "text/csv")
			c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, ue.Nickname))
			_, err = c.Response().Write([]byte(*v.TextData))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, nil)
			}
		}
	}
	return nil
}

func ExportRunCsvRequest2(ctx context.Context, ou org_users.OrgUser, id int) ([]artemis_entities.UserEntity, error) {
	ojsRuns, err := artemis_orchestrations.AdminSelectAiSystemOrchestrations(ctx, ou, id)
	if err != nil {
		log.Err(err).Msg("failed to get runs")
		return nil, err
	}
	if len(ojsRuns) < 1 {
		return nil, nil // todo err
	}
	ojr := ojsRuns[0]
	// 7138983863666903883
	//ojr := artemis_orchestrations.OrchestrationsAnalysis{}
	//ojr.OrchestrationName = "test-wf"
	log.Info().Interface("oj.GroupName", ojr.GroupName).Msg("ExportRunCsvRequest")
	var ueSlice []artemis_entities.UserEntity

	ue := artemis_entities.UserEntity{
		Nickname: ojr.OrchestrationName,
		Platform: "csv-exports",
	}
	_, err = ai_platform_service_orchestrations.S3WfRunExport(ctx, ou, ojr.OrchestrationName, &ue)
	if err != nil {
		return nil, err
	}
	ue.Nickname = ojr.GroupName
	p := &filepaths.Path{
		DirOut: fmt.Sprintf("/debug/runs/%s", ojr.OrchestrationName),
		FnOut:  fmt.Sprintf("wf_entrypoint.json"),
	}
	b, err := ai_platform_service_orchestrations.S3WfRunDownloadWithPath(ctx, p)
	if err != nil {
		log.Warn().Err(err).Msg("ExportRunCsvRequest2: S3WfRunDownloadWithPath")
		return nil, err
	}

	if b.Bytes() != nil {
		flowIn := &ExecFlowsActionsRequest{}
		err = json.Unmarshal(b.Bytes(), flowIn)
		if err != nil {
			log.Warn().Err(err).Msg("ExportRunCsvRequest2: Unmarshal")
			return nil, err
		}
		ueInputs := artemis_entities.UserEntity{
			Nickname: fmt.Sprintf("%s_inputs", ojr.GroupName),
			Platform: "csv-exports",
			MdSlice: []artemis_entities.UserEntityMetadata{
				{
					TextData: aws.String(flowIn.ContactsCsvStr),
				},
			},
		}
		ueSlice = append(ueSlice, ueInputs)
		uePrompts := artemis_entities.UserEntity{
			Nickname: fmt.Sprintf("%s_prompts", ojr.GroupName),
			Platform: "csv-exports",
			MdSlice: []artemis_entities.UserEntityMetadata{
				{
					TextData: aws.String(flowIn.PromptsCsvStr),
				},
			},
		}
		ueSlice = append(ueSlice, uePrompts)
	}
	ueSlice = append(ueSlice, ue)
	return ueSlice, nil
}