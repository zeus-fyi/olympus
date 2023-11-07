package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
)

func ProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Provision(c)
}

type ProvisionRequest struct {
	hestia_quicknode.ProvisionRequest
	hestia_quicknode.QuickNodeUserInfo `json:"verified"`
}

const (
	TestPlan        = "test"
	FreePlan        = "free"
	DiscoverPlan    = "discover"
	LitePlan        = "lite"
	Standard        = "standard"
	PerformancePlan = "performance"
)

func (r *ProvisionRequest) Provision(c echo.Context) error {
	if c.Get(QuickNodeIDHeader) == nil {
		key, err := auth.VerifyQuickNodeToken(context.Background(), r.QuickNodeID)
		if err != nil {
			log.Err(err).Msg("InitV1Routes QuickNode user not found: creating new org")
			err = nil
		}
		ou := org_users.NewOrgUserWithID(key.OrgID, 0)
		c.Set("orgUser", ou)
		c.Set("verified", key.IsVerified())
	}

	ouc := c.Get("orgUser")
	ou, ok := ouc.(org_users.OrgUser)
	if !ok {
		key, err := auth.VerifyQuickNodeToken(context.Background(), r.QuickNodeID)
		if err != nil {
			log.Err(err).Msg("InitV1Routes QuickNode user not found: creating new org")
			err = nil
		}
		ou = org_users.NewOrgUserWithID(key.OrgID, 0)
		c.Set("orgUser", ou)
		c.Set("verified", key.IsVerified())
	}
	pr := r.ProvisionRequest
	r.Verified = false
	val, ok := c.Get("verified").(bool)
	if ok {
		r.Verified = val
	}
	isTestReq, ok := c.Get("isTest").(bool)
	if ok {
		r.IsTest = isTestReq
	} else {
		r.IsTest = false
	}
	switch pr.Plan {
	case LitePlan, Standard, PerformancePlan:
	case TestPlan:
		if !r.IsTest {
			return c.JSON(http.StatusBadRequest, QuickNodeResponse{
				Error: "error: plan not supported",
			})
		}
	default:
		return c.JSON(http.StatusBadRequest, QuickNodeResponse{
			Error: "error: plan not supported",
		})
	}

	if pr.QuickNodeID == "" {
		return c.JSON(http.StatusBadRequest, QuickNodeResponse{
			Error: "error: quicknode id not provided",
		})
	}

	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnProvisionWorkflow(context.Background(), ou, pr, r.QuickNodeUserInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Error: "error: internal error: failed to provision quicknode service",
			})
	}
	return c.JSON(http.StatusOK, ProvisionResponse{
		AccessURL:    "https://iris.zeus.fyi/v1/router",
		DashboardURL: "https://cloud.zeus.fyi/quicknode/dashboard",
		Status:       "success",
	})
}

func TestProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ProvisionTest(c)
}

func (r *ProvisionRequest) ProvisionTest(c echo.Context) error {
	ouc := c.Get("orgUser")
	ou, ok := ouc.(org_users.OrgUser)
	if !ok {
		ou = org_users.OrgUser{}
		key, err := auth.VerifyQuickNodeToken(context.Background(), r.QuickNodeID)
		if err != nil {
			log.Err(err).Msg("InitV1Routes QuickNode user not found: creating new org")
			err = nil
		}
		ou = org_users.NewOrgUserWithID(key.OrgID, 0)
		c.Set("orgUser", ou)
		c.Set("verified", key.IsVerified())
	}

	pr := r.ProvisionRequest
	r.Verified = false
	val, ok := c.Get("verified").(bool)
	if ok {
		r.Verified = val
	}
	isTestReq, ok := c.Get("isTest").(bool)
	if ok {
		r.IsTest = isTestReq
	} else {
		r.IsTest = false
	}
	pr.Plan = "lite"
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnProvisionWorkflow(context.Background(), ou, pr, r.QuickNodeUserInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Error: "error: internal error: failed to provision quicknode service",
			})
	}

	return c.JSON(http.StatusOK, ProvisionResponse{
		AccessURL:    "https://iris.zeus.fyi/v1/router",
		DashboardURL: "http://localhost:3000/quicknode/dashboard",
		Status:       "success",
	})
}

type ProvisionResponse struct {
	Status       string `json:"status"`
	DashboardURL string `json:"dashboard-url"`
	AccessURL    string `json:"access-url"`
}

func UpdateProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.UpdateProvision(c)
}

func (r *ProvisionRequest) UpdateProvision(c echo.Context) error {
	ouc := c.Get("orgUser")
	ou, ok := ouc.(org_users.OrgUser)
	if !ok {
		key, err := auth.VerifyQuickNodeToken(context.Background(), r.QuickNodeID)
		if err != nil {
			log.Err(err).Msg("InitV1Routes QuickNode user not found: creating new org")
			err = nil
		}
		ou = org_users.NewOrgUserWithID(key.OrgID, 0)
		c.Set("orgUser", ou)
		c.Set("verified", key.IsVerified())
	}
	r.Verified = false
	val, ok := c.Get("verified").(bool)
	if ok {
		r.Verified = val
	}
	isTestReq, ok := c.Get("isTest").(bool)
	if ok {
		r.IsTest = isTestReq
	} else {
		r.IsTest = false
	}
	pr := r.ProvisionRequest
	switch pr.Plan {
	case LitePlan, Standard, PerformancePlan, DiscoverPlan, FreePlan:
	case TestPlan:
		if !r.IsTest {
			return c.JSON(http.StatusBadRequest, QuickNodeResponse{
				Error: "error: plan not supported",
			})
		}
	default:
		return c.JSON(http.StatusBadRequest, QuickNodeResponse{
			Error: "error: plan not supported",
		})
	}

	if pr.QuickNodeID == "" {
		return c.JSON(http.StatusBadRequest, QuickNodeResponse{
			Error: "error: quicknode id not provided",
		})
	}

	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnUpdateProvisionWorkflow(context.Background(), ou, pr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Status: "error",
			})
	}
	return c.JSON(http.StatusOK, ProvisionResponse{
		AccessURL:    "https://iris.zeus.fyi/v1/router",
		DashboardURL: "https://cloud.zeus.fyi/quicknode/dashboard",
		Status:       "success",
	})
}
