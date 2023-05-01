package hestia_infracost

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type InfraCostClient struct {
	*resty.Client
}

func InitInfraCostClient(ctx context.Context, apiKey string) InfraCostClient {
	r := resty.New()
	r.SetHeader("X-Api-Key", apiKey)
	return InfraCostClient{
		r,
	}
}

type ProductFilter struct {
	VendorName       string             `json:"vendorName" graphql:"vendorName"`
	Service          string             `json:"service" graphql:"service"`
	ProductFamily    string             `json:"productFamily" graphql:"productFamily"`
	Region           string             `json:"region" graphql:"region"`
	SKU              string             `json:"sku,omitempty" graphql:"sku"`
	AttributeFilters []*AttributeFilter `json:"attributeFilters" graphql:"attributeFilters"`
}

type Attribute struct {
	Key   string `json:"key" graphql:"key"`
	Value string `json:"value" graphql:"value"`
}

type AttributeFilter struct {
	Key        string `json:"key" graphql:"key"`
	Value      string `json:"value" graphql:"value"`
	ValueRegex string `json:"value_regex,omitempty" graphql:"value_regex"`
}

type PriceFilter struct {
	PurchaseOption     string `json:"purchaseOption" graphql:"purchaseOption"`
	Unit               string `json:"unit" graphql:"unit"`
	Description        string `json:"description" graphql:"description"`
	DescriptionRegex   string `json:"description_regex" graphql:"description_regex"`
	StartUsageAmount   string `json:"startUsageAmount" graphql:"startUsageAmount"`
	EndUsageAmount     string `json:"endUsageAmount" graphql:"endUsageAmount"`
	TermLength         string `json:"termLength" graphql:"termLength"`
	TermPurchaseOption string `json:"termPurchaseOption" graphql:"termPurchaseOption"`
	TermOfferingClass  string `json:"termOfferingClass" graphql:"termOfferingClass"`
}

type ProductsRequest struct {
	Query     string        `json:"query"`
	Variables ProductsInput `json:"variables"`
}

type ProductsInput struct {
	Filter ProductFilter `json:"filter"`
}

func (i *InfraCostClient) GetCost(ctx context.Context, p ProductFilter) error {
	requestURL := "https://pricing.api.infracost.io/graphql"

	productsReq := ProductsRequest{
		Query: `
            query GetProducts($filter: ProductFilter!) { 
                products(filter: $filter) { 
                    prices(filter: { purchaseOption: "on_demand" }) { 
                        USD 
                    } 
                } 
            }
        `,
		Variables: ProductsInput{
			Filter: p,
		},
	}
	// Send the request using the Post() method
	resp, err := i.R().
		SetHeader("Content-Type", "application/json").
		SetBody(productsReq).
		Post(requestURL)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		return err
	}

	fmt.Println(resp.String())
	// Access the result as needed
	//fmt.Printf("USD price: %f\n", result.Data.Products[0].Prices[0].USD)
	return err
}
