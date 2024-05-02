package zeus_v1_ai

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	ai_platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/zeus/ai/orchestrations"
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
