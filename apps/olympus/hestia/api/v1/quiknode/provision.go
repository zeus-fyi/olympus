package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	Verified bool `json:"verified"`
}

const (
	TrialPlan       = "trial"
	LitePlan        = "lite"
	Standard        = "standard"
	PerformancePlan = "performance"
)

func (r *ProvisionRequest) Provision(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	pr := r.ProvisionRequest
	r.Verified = false
	val, ok := c.Get("verified").(bool)
	if ok {
		r.Verified = val
	}
	switch pr.Plan {
	case LitePlan, Standard, PerformancePlan:
	default:
		return c.JSON(http.StatusBadRequest, QuickNodeResponse{
			Error: "error: plan not supported",
		})
	}
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnProvisionWorkflow(context.Background(), pr, ou, r.Verified)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Error: "error: internal error: failed to provision quicknode service",
			})
	}
	return c.JSON(http.StatusOK, ProvisionResponse{
		AccessURL:    "https://cloud.zeus.fyi/quicknode/access",
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
	return c.JSON(http.StatusOK, ProvisionResponse{
		AccessURL:    "http://localhost:9002/quicknode/access",
		DashboardURL: "http://localhost:9002/quicknode/dashboard",
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
	ou := c.Get("orgUser").(org_users.OrgUser)
	pr := r.ProvisionRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnUpdateProvisionWorkflow(context.Background(), pr, ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Status: "error",
			})
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}
