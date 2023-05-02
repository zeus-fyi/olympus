package hestia_gcp

type Category struct {
	ServiceDisplayName string `json:"serviceDisplayName"`
	ResourceFamily     string `json:"resourceFamily"`
	ResourceGroup      string `json:"resourceGroup"`
	UsageType          string `json:"usageType"`
}

type Money struct {
	CurrencyCode string `json:"currencyCode"`
	Units        string `json:"units"`
	Nanos        int    `json:"nanos"`
}

type TierRate struct {
	StartUsageAmount float64 `json:"startUsageAmount"`
	UnitPrice        Money   `json:"unitPrice"`
}

type PricingExpression struct {
	UsageUnit                string     `json:"usageUnit"`
	UsageUnitDescription     string     `json:"usageUnitDescription"`
	BaseUnit                 string     `json:"baseUnit"`
	BaseUnitDescription      string     `json:"baseUnitDescription"`
	BaseUnitConversionFactor float64    `json:"baseUnitConversionFactor"`
	DisplayQuantity          float64    `json:"displayQuantity"`
	TieredRates              []TierRate `json:"tieredRates"`
}

type PricingInfo struct {
	EffectiveTime          string            `json:"effectiveTime"`
	Summary                string            `json:"summary"`
	PricingExpression      PricingExpression `json:"pricingExpression"`
	CurrencyConversionRate float64           `json:"currencyConversionRate"`
}

type GeoTaxonomy struct {
	Type    string   `json:"type"`
	Regions []string `json:"regions"`
}

type Sku struct {
	Name                string        `json:"name"`
	SkuId               string        `json:"skuId"`
	Description         string        `json:"description"`
	Category            Category      `json:"category"`
	ServiceRegions      []string      `json:"serviceRegions"`
	PricingInfo         []PricingInfo `json:"pricingInfo"`
	ServiceProviderName string        `json:"serviceProviderName"`
	GeoTaxonomy         GeoTaxonomy   `json:"geoTaxonomy"`
}

type SkusResponse struct {
	Skus          []Sku  `json:"skus"`
	NextPageToken string `json:"nextPageToken"`
}
