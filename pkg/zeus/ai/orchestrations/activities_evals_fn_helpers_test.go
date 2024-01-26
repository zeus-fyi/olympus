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
	}
	//rerr := act.EvalModelScoredJsonOutput(ctx, evalFns)
	//t.Require().Nil(rerr)

	for _, evalFnWithMetrics := range evalFns {
		for _, fi := range evalFnWithMetrics.Schemas {
			for fieldInd, f := range fi.Fields {
				fmt.Println(f.FieldName, f.DataType, f.FieldValue)
				switch f.DataType {
				case "string":
					t.Require().Equal("test", *fi.Fields[fieldInd].StringValue)
					t.Require().True(fi.Fields[fieldInd].IsValidated)
				case "number":
					t.Require().Equal(1.0, *fi.Fields[fieldInd].NumberValue)
					t.Require().True(fi.Fields[fieldInd].IsValidated)
				}
				t.Require().NotNil(f.EvalMetrics)
				t.Require().NotEmpty(f.EvalMetrics)

				for _, evm := range f.EvalMetrics {
					t.Require().NotEmpty(evm)
					t.Require().NotNil(evm)
					t.Require().NotNil(evm.EvalMetricResult)
					t.Require().NotNil(evm.EvalMetricResult.EvalResultOutcomeBool)
					fmt.Println(f.FieldName, *evm.EvalMetricResult.EvalResultOutcomeBool)
				}
			}
		}
	}
}

func (t *ZeusWorkerTestSuite) TestEvalModelScoredJsonOutput() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	act := NewZeusAiPlatformActivities()
	evalFns, err := act.EvalLookup(ctx, ou, 1704066747085827000)
	t.Require().Nil(err)
	t.Require().NotEmpty(evalFns)

	t.Require().Nil(err)
}
