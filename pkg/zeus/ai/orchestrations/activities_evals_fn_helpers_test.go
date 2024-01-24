package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (t *ZeusWorkerTestSuite) TestTransformJSONToEvalScoredMetrics() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFns, err := act.EvalLookup(ctx, ou, 1704066747085827000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFns)

	for _, evalFnWithMetrics := range evalFns {
		for _, fi := range evalFnWithMetrics.Schemas {
			for fieldInd, f := range fi.Fields {
				fmt.Println(f.FieldName, f.DataType, f.FieldValue)
				switch f.DataType {
				case "string":
					fi.Fields[fieldInd].StringValue = aws.String("test")
					fi.Fields[fieldInd].IsValidated = true
				case "number":
					fi.Fields[fieldInd].NumberValue = aws.Float64(1.0)
					fi.Fields[fieldInd].IsValidated = true
				}
			}
		}
		fmt.Println(evalFnWithMetrics.EvalType)
		rerr := act.EvalModelScoredJsonOutput(ctx, &evalFnWithMetrics)
		t.Require().Nil(rerr)
	}
}
