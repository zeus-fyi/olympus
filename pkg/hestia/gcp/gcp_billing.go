package hestia_gcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
)

const (
	GcpServicesURL = "https://cloudbilling.googleapis.com/v1/services"
)

type ServicesResponse struct {
	Services      hestia_autogen_bases.GcpServicesSlice `json:"services"`
	NextPageToken string                                `json:"nextPageToken"`
}

func (g *GcpClient) ListServices(ctx context.Context) (ServicesResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("pageToken", "")
	totalServices := ServicesResponse{}
	for {
		requestURL := fmt.Sprintf(GcpServicesURL)
		requestURL = fmt.Sprintf("%s?%s", requestURL, queryParams.Encode())
		// Execute the request

		var tmp ServicesResponse
		resp, err := g.R().SetResult(&tmp).Get(requestURL)
		if err != nil {
			fmt.Printf("Error executing request: %v\n", err)
			return totalServices, err
		}
		// Check for non-2xx status codes
		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			fmt.Printf("Error: API responded with status code %d\n", resp.StatusCode())
			return totalServices, err
		}
		totalServices.Services = append(totalServices.Services, tmp.Services...)
		if tmp.NextPageToken == "" {
			return totalServices, nil
		}
		queryParams.Set("pageToken", tmp.NextPageToken)
	}
}

func (g *GcpClient) ListServiceSKUs(ctx context.Context, parentService, serviceID string) error {
	queryParams := url.Values{}
	queryParams.Set("pageToken", "")
	totalSKUs := SkusResponse{}
	for {
		requestURL := fmt.Sprintf("https://cloudbilling.googleapis.com/v1/%s/skus", parentService)
		requestURL = fmt.Sprintf("%s?%s", requestURL, queryParams.Encode())
		// Execute the request

		var tmp SkusResponse
		resp, err := g.R().SetResult(&tmp).Get(requestURL)
		if err != nil {
			fmt.Printf("Error executing request: %v\n", err)
			return err
		}
		// Check for non-2xx status codes
		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			fmt.Printf("Error: API responded with status code %d\n", resp.StatusCode())
			return err
		}
		totalSKUs.Skus = append(totalSKUs.Skus, tmp.Skus...)
		if tmp.NextPageToken == "" {
			for _, sku := range totalSKUs.Skus {
				insertSku, ierr := convertSkuToGcpServicesSkus(serviceID, sku)
				if ierr != nil {
					log.Ctx(ctx).Error().Err(ierr).Msg("Error convertSkuToGcpServicesSkus")
					return err
				}
				err = hestia_compute_resources.InsertGcpServicesSKU(ctx, insertSku)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("Error inserting GCP SKU")
					return err
				}
			}
			return nil
		}
		queryParams.Set("pageToken", tmp.NextPageToken)
	}
}

func convertSkuToGcpServicesSkus(serviceID string, sku Sku) (autogen_bases.GcpServicesSkus, error) {
	serviceRegions, err := json.Marshal(sku.ServiceRegions)
	if err != nil {
		return autogen_bases.GcpServicesSkus{}, err
	}

	pricingInfo, err := json.Marshal(sku.PricingInfo)
	if err != nil {
		return autogen_bases.GcpServicesSkus{}, err
	}

	geoTaxonomy, err := json.Marshal(sku.GeoTaxonomy)
	if err != nil {
		return autogen_bases.GcpServicesSkus{}, err
	}

	return autogen_bases.GcpServicesSkus{
		UsageType:           sql.NullString{String: sku.Category.UsageType, Valid: sku.Category.UsageType != ""},
		ServiceRegions:      sql.NullString{String: string(serviceRegions), Valid: len(sku.ServiceRegions) > 0},
		PricingInfo:         sql.NullString{String: string(pricingInfo), Valid: len(sku.PricingInfo) > 0},
		GeoTaxonomy:         sql.NullString{String: string(geoTaxonomy), Valid: sku.GeoTaxonomy.Type != ""},
		ServiceID:           serviceID,
		ServiceDisplayName:  sql.NullString{String: sku.Category.ServiceDisplayName, Valid: sku.Category.ServiceDisplayName != ""},
		ResourceFamily:      sql.NullString{String: sku.Category.ResourceFamily, Valid: sku.Category.ResourceFamily != ""},
		ResourceGroup:       sql.NullString{String: sku.Category.ResourceGroup, Valid: sku.Category.ResourceGroup != ""},
		ServiceProviderName: sql.NullString{String: sku.ServiceProviderName, Valid: sku.ServiceProviderName != ""},
		Name:                sku.Name,
		SkuID:               sku.SkuId,
		Description:         sql.NullString{String: sku.Description, Valid: sku.Description != ""},
	}, nil
}
