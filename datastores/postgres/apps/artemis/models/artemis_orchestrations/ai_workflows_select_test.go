package artemis_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplate() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow2",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup2",
		FundamentalPeriodTimeUnit: "days",
	}

	res, err := SelectWorkflowTemplate(ctx, ou, newTemplate.WorkflowName)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)

	md := MapDependencies(res)
	fmt.Println("\nAnalysis")

	if newTemplate.WorkflowName == "Example Workflow1" {
		fmt.Println("newTemplate.WorkflowName", newTemplate.WorkflowName)
		s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701653245709972992])
		s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701667813254964224])
		s.Require().Equal(true, md.AnalysisRetrievals[1701657795016150016][1701667784112279040])

		for k, v := range md.AnalysisRetrievals {
			fmt.Println(k, v)
		}
		fmt.Println("\nAgg")
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
		for k, v := range md.AggregateAnalysis {
			fmt.Println(k, v)
		}
	}
	if newTemplate.WorkflowName == "Example Workflow2" {
		s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701653245709972992])
		s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701667813254964224])
		s.Require().Equal(true, md.AnalysisRetrievals[1701657795016150016][1701667784112279040])

		for k, v := range md.AnalysisRetrievals {
			fmt.Println(k, v)
		}
		fmt.Println("\nAgg")
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
		s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
		for k, v := range md.AggregateAnalysis {
			fmt.Println(k, v)
		}
	}
}

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplatesP() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	newTemplate := WorkflowTemplate{
		WorkflowName:  "wf-test-health",
		WorkflowGroup: "wf-test-health",
	}

	res1, err := SelectWorkflowTemplate(ctx, ou, newTemplate.WorkflowName)
	s.Require().Nil(err)
	s.Require().NotEmpty(res1)
	md := MapDependencies(res1)
	s.Require().NotEmpty(md.AnalysisRetrievals)
	res, err := SelectWorkflowTemplates(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
}

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplates() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectWorkflowTemplates(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)

	for _, newTemplate := range res.WorkflowTemplatesMap {

		md := MapDependenciesGrouped(newTemplate)

		if newTemplate.WorkflowName == "Example Workflow1" {
			fmt.Println("newTemplate.WorkflowName", newTemplate.WorkflowName)

			s.Require().Equal(true, md.AnalysisRetrievals[1701657795016150016][1701667784112279040])
			s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701653245709972992])
			s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701667813254964224])
			for k, v := range md.AnalysisRetrievals {
				fmt.Println(k, v)
			}
			fmt.Println("\nAgg")
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
			for k, v := range md.AggregateAnalysis {
				fmt.Println(k, v)
			}
		}
		if newTemplate.WorkflowName == "Example Workflow2" {
			s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701653245709972992])
			s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701667813254964224])
			s.Require().Equal(true, md.AnalysisRetrievals[1701657795016150016][1701667784112279040])

			for k, v := range md.AnalysisRetrievals {
				fmt.Println(k, v)
			}
			fmt.Println("\nAgg")
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657822027992064])
			s.Require().Equal(true, md.AggregateAnalysis[1701657830780669952][1701657795016150016])
			for k, v := range md.AggregateAnalysis {
				fmt.Println(k, v)
			}
		}

		if newTemplate.WorkflowName == "Example Workflow4" {
			s.Require().Equal(true, md.AnalysisRetrievals[1701657822027992064][1701667813254964224])

			s.Require().Equal(1, len(md.AnalysisRetrievals))
			for k, v := range md.AnalysisRetrievals {
				fmt.Println(k, v)
			}
			fmt.Println("\nAgg")
			s.Require().Equal(0, len(md.AggregateAnalysis))

			for k, v := range md.AggregateAnalysis {
				fmt.Println(k, v)
			}
		}

		if newTemplate.WorkflowName == "Example Workflow With EvalFns" {
			count := 0
			evalFns := newTemplate.AnalysisEvalFns[1701657822027992064]
			s.Assert().NotEmpty(evalFns)
			for _, ef := range evalFns {
				s.Require().Equal(1703624059411640000, ef.EvalID)
				count += 1
			}

			s.Assert().Equal(1, count)
		}

		if newTemplate.WorkflowName == "Test Workflow With Agg EvalFns 6" {
			count := 0
			aggEvalFns := newTemplate.AggEvalFns[1701657830780669952]
			s.Assert().NotEmpty(aggEvalFns)

			for _, ef := range aggEvalFns {
				s.Require().Equal(1703624059411640000, ef.EvalID)
				count += 1
			}

			evalFn := newTemplate.AggAnalysisEvalFns[1701657830780669952][1701657822027992064]
			s.Assert().Equal(1703624059411640000, evalFn.EvalID)

			evalFn2 := newTemplate.AggAnalysisEvalFns[1701657830780669952][1701657795016150016]
			s.Assert().Equal(1702959527789976000, evalFn2.EvalID)

			s.Require().Equal(1, count)
		}
	}
}

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplateSingle() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	name := "Test Workflow With Agg EvalFns 6"
	res, err := SelectWorkflowTemplateByName(ctx, ou, name)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)

	for _, newTemplate := range res.WorkflowTemplatesMap {

		aggEvalFns := newTemplate.AggEvalFns[1701657830780669952]
		s.Assert().NotEmpty(aggEvalFns)
		count := 0
		for _, ef := range aggEvalFns {
			s.Require().Equal(1703624059411640000, ef.EvalID)
			count += 1
		}

		evalFn := newTemplate.AggAnalysisEvalFns[1701657830780669952][1701657822027992064]
		s.Assert().Equal(1703624059411640000, evalFn.EvalID)

		evalFn2 := newTemplate.AggAnalysisEvalFns[1701657830780669952][1701657795016150016]
		s.Assert().Equal(1702959527789976000, evalFn2.EvalID)

		s.Require().Equal(1, count)
	}
}

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplatesWithEvalFns() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	res, err := SelectWorkflowTemplates(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
}
