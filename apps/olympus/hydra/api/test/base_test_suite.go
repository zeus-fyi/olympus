package hydra_base_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var Kns zeus_common_types.CloudCtxNs

var TestOrgUser = org_users.OrgUser{
	OrgUsers: autogen_bases.OrgUsers{OrgID: 1667452524363177528, UserID: 1667452524356256466},
}

var TestTopologyID = 7140168037686545724

type HydraBaseTestSuite struct {
	E  *echo.Echo
	Eg *echo.Group
	autok8s_core.K8TestSuite
	D        test_suites.DatastoresTestSuite
	Ts       chronos.Chronos
	Endpoint string

	zeus_client.ZeusClient
}

func (t *HydraBaseTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.E = echo.New()

	t.ZeusClient = zeus_client.NewZeusClient("http://localhost:9000", "")
}

func (t *HydraBaseTestSuite) PostRequest(ctx context.Context, endpoint string, payload any) ([]byte, error) {
	t.PrintReqJson(payload)

	resp, err := t.R().
		SetBody(payload).
		Post(endpoint)

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Ctx(ctx).Err(err).Msg("HydraBaseTestSuite: PostRequest")
		if resp.StatusCode() == http.StatusBadRequest {
			err = errors.New("bad request")
		}
		return nil, err
	}
	t.PrintRespJson(resp.Body())
	return resp.Body(), err
}
