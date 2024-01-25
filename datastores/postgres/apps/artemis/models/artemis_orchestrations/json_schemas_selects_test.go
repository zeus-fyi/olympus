package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectJsonSchemaByOrg() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	// get internal assignments
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	js, err := SelectJsonSchemaByOrg(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotNil(js)

	count := 0
	for _, j := range js.Slice {
		if j.SchemaID != 1706143175198902000 {
			continue
		}
		for _, f := range j.Fields {
			switch f.FieldName {
			case "msg_id":
				count += 10
				s.Require().Equal(1706143175298912935, f.FieldID)
			case "msg_score":
				count += 1
				s.Require().Equal(1706143175335893834, f.FieldID)
			}
		}
	}
	s.Require().Equal(11, count)
}
