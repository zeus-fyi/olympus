package hestia_eks_aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"github.com/rs/zerolog/log"
)

type AwsPricing struct {
	*pricing.Client
}

func InitPricingClient(ctx context.Context, accessCred EksCredentials) (AwsPricing, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.AccessSecret, "")
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
					Value: aws.String("t3"),
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

func (a *AwsPricing) GetEC2Product(ctx context.Context, region, instanceType string) (AWSPrice, error) {
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
			log.Ctx(ctx).Err(err)
			return AWSPrice{}, err
		}

		for _, prod := range pa.PriceList {
			ec2Prod := AWSPrice{}
			err = json.Unmarshal([]byte(prod), &ec2Prod)
			if err != nil {
				return ec2Prod, err
			}
			return ec2Prod, nil
		}
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
