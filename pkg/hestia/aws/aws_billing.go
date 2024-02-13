package hestia_eks_aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"github.com/rs/zerolog/log"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsPricing struct {
	*pricing.Client
}

func InitPricingClient(ctx context.Context, accessCred aegis_aws_auth.AuthAWS) (AwsPricing, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.SecretKey, "")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(accessCred.Region),
	)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return AwsPricing{}, err
	}
	return AwsPricing{pricing.NewFromConfig(cfg)}, nil
}

func (a *AwsPricing) GetAllProducts(ctx context.Context, region string) error {
	for {
		pi := &pricing.GetProductsInput{
			ServiceCode: aws.String("AmazonEC2"),
			Filters: []types.Filter{
				{
					Field: aws.String("regionCode"),
					Type:  "TERM_MATCH",
					Value: aws.String(region),
				},
				{
					Field: aws.String("ServiceCode"),
					Type:  "TERM_MATCH",
					Value: aws.String("AmazonEC2"),
				},
				{
					Field: aws.String("marketoption"),
					Type:  "TERM_MATCH",
					Value: aws.String("OnDemand"),
				},
				{
					Field: aws.String("instanceType"),
					Type:  "TERM_MATCH",
					Value: aws.String("m6a"),
				},
			},
			NextToken: nil,
		}
		pa, err := a.GetProducts(ctx, pi)
		if err != nil {
			log.Ctx(ctx).Err(err)
			return err
		}

		for _, prod := range pa.PriceList {
			ec2Prod := AWSPrice{}
			err = json.Unmarshal([]byte(prod), &ec2Prod)
			if err != nil {
				return err
			}
			//fmt.Println(ec2Prod.Product.Attributes["instanceType"])
		}
		//fmt.Println(pa.NextToken)
		pi.NextToken = pa.NextToken
		if pa.NextToken == nil {
			return nil
		}
	}
}

func (a *AwsPricing) GetEC2Product(ctx context.Context, region, instanceType string) ([]AWSPrice, error) {
	for {
		pi := &pricing.GetProductsInput{
			ServiceCode: aws.String("AmazonEC2"),
			Filters: []types.Filter{
				{
					Field: aws.String("regionCode"),
					Type:  "TERM_MATCH",
					Value: aws.String(region),
				},
				{
					Field: aws.String("ServiceCode"),
					Type:  "TERM_MATCH",
					Value: aws.String("AmazonEC2"),
				},
				{
					Field: aws.String("marketoption"),
					Type:  "TERM_MATCH",
					Value: aws.String("OnDemand"),
				},
				{
					Field: aws.String("instanceType"),
					Type:  "TERM_MATCH",
					Value: aws.String(instanceType),
				},
			},
			NextToken: nil,
		}
		pa, err := a.GetProducts(ctx, pi)
		if err != nil {
			log.Err(err)
			return []AWSPrice{}, err
		}

		var prices []AWSPrice
		for _, prod := range pa.PriceList {
			ec2Prod := AWSPrice{}
			err = json.Unmarshal([]byte(prod), &ec2Prod)
			if err != nil {
				return prices, err
			}
			prices = append(prices, ec2Prod)
		}
		return prices, nil
	}
}

type AWSPrice struct {
	Product         Product `json:"product"`
	ServiceCode     string  `json:"serviceCode"`
	Terms           Terms   `json:"terms"`
	Version         string  `json:"version"`
	PublicationDate string  `json:"publicationDate"`
}

type Product struct {
	ProductFamily string            `json:"productFamily"`
	Attributes    map[string]string `json:"attributes"`
	Sku           string            `json:"sku"`
}

type Terms struct {
	OnDemand map[string]OnDemandTerm `json:"OnDemand"`
}

type OnDemandTerm struct {
	PriceDimensions map[string]PriceDimension `json:"priceDimensions"`
	Sku             string                    `json:"sku"`
	EffectiveDate   string                    `json:"effectiveDate"`
	OfferTermCode   string                    `json:"offerTermCode"`
	TermAttributes  map[string]string         `json:"termAttributes"`
}

type PriceDimension struct {
	Unit         string            `json:"unit"`
	EndRange     string            `json:"endRange"`
	Description  string            `json:"description"`
	AppliesTo    []interface{}     `json:"appliesTo"`
	RateCode     string            `json:"rateCode"`
	BeginRange   string            `json:"beginRange"`
	PricePerUnit map[string]string `json:"pricePerUnit"`
}

func (ap *AWSPrice) GetPricePerUnitUSD() (float64, string, error) {
	for _, v := range ap.Terms.OnDemand {
		for _, pd := range v.PriceDimensions {
			if val, ok := pd.PricePerUnit["USD"]; ok {
				price, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return 0, "", err
				}
				return price, pd.Unit, nil
			}
		}
	}
	return 0, "", fmt.Errorf("no USD price per unit found")
}

func (ap *AWSPrice) GetDescription() string {
	for _, v := range ap.Terms.OnDemand {
		for _, pd := range v.PriceDimensions {
			return pd.Description
		}
	}
	return ""
}

func (ap *AWSPrice) GetVCpus() string {
	return ap.Product.Attributes["vcpu"]
}

func (ap *AWSPrice) GetMemoryAndUnits() (string, string) {
	memory := strings.Split(ap.Product.Attributes["memory"], " ")

	if len(memory) < 2 {
		return "", ""
	}
	return memory[0], memory[1]
}
