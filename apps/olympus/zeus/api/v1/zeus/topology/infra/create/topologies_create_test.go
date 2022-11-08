package create_infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/pretty"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyCreateActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
	c conversions_test.ConversionsTestSuite
	h hestia_test.BaseHestiaTestSuite
}

func (t *TopologyCreateActionRequestTestSuite) TestUpload() {
	t.InitLocalConfigs()
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	name := fmt.Sprintf("random_%d", t.Ts.UnixTimeStampNow())
	c := charts.Chart{}
	c.ChartName = "test_api"
	c.ChartVersion = fmt.Sprintf("test_api_v%d", t.Ts.UnixTimeStampNow())
	oid, uid := t.h.NewTestOrgAndUser()
	//orgUser := org_users.NewOrgUserWithID(oid, uid)
	fmt.Printf("orgID: %d\n userID %d\n", oid, uid)

	createRequest := TopologyCreateRequest{
		TopologyName:     name,
		ChartName:        c.ChartName,
		ChartDescription: "describes chart",
		Version:          fmt.Sprintf("v0.0.%d", +t.Ts.UnixTimeStampNow()),
	}
	topologyActionRequestPayload, err := json.Marshal(createRequest)
	t.Assert().Nil(err)

	fmt.Println("action request json")
	requestJSON := pretty.Pretty(topologyActionRequestPayload)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))

	t.E.POST("/infra", nil)
	client := resty.New()
	resp, err := client.R().
		SetFile("chart", "./zeus.tar.gz").
		Post("http://localhost:9010/infra")

	t.Require().Nil(err)
	t.Assert().Equal(http.StatusOK, resp.StatusCode())

	result := pretty.Pretty(resp.Body())
	result = pretty.Color(result, pretty.TerminalStyle)
	fmt.Println(string(result))
}

func TestTopologyCreateActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyCreateActionRequestTestSuite))
}
