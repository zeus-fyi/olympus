package price_quoter

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SwapQuote struct {
	ChainID              int      `json:"chainId"`
	Price                string   `json:"price"`
	GuaranteedPrice      string   `json:"guaranteedPrice"`
	EstimatedPriceImpact string   `json:"estimatedPriceImpact"`
	To                   string   `json:"to"`
	Data                 string   `json:"data"`
	Value                string   `json:"value"`
	Gas                  string   `json:"gas"`
	EstimatedGas         string   `json:"estimatedGas"`
	GasPrice             string   `json:"gasPrice"`
	ProtocolFee          string   `json:"protocolFee"`
	MinimumProtocolFee   string   `json:"minimumProtocolFee"`
	BuyTokenAddress      string   `json:"buyTokenAddress"`
	SellTokenAddress     string   `json:"sellTokenAddress"`
	BuyAmount            string   `json:"buyAmount"`
	SellAmount           string   `json:"sellAmount"`
	Sources              []Source `json:"sources"`
	Orders               []Order  `json:"orders"`
	AllowanceTarget      string   `json:"allowanceTarget"`
	DecodedUniqueId      string   `json:"decodedUniqueId"`
	SellTokenToEthRate   string   `json:"sellTokenToEthRate"`
	BuyTokenToEthRate    string   `json:"buyTokenToEthRate"`
	AuxiliaryChainData   struct{} `json:"auxiliaryChainData"`
	ExpectedSlippage     string   `json:"expectedSlippage"`
}

type Source struct {
	Name       string `json:"name"`
	Proportion string `json:"proportion"`
}

type Order struct {
	Type        int      `json:"type"`
	Source      string   `json:"source"`
	MakerToken  string   `json:"makerToken"`
	TakerToken  string   `json:"takerToken"`
	MakerAmount string   `json:"makerAmount"`
	TakerAmount string   `json:"takerAmount"`
	FillData    FillData `json:"fillData"`
	Fill        Fill     `json:"fill"`
}

type FillData struct {
	TokenAddressPath []string `json:"tokenAddressPath"`
	Router           string   `json:"router"`
}

type Fill struct {
	Input          string `json:"input"`
	Output         string `json:"output"`
	AdjustedOutput string `json:"adjustedOutput"`
	Gas            int    `json:"gas"`
}

type Client struct {
	http *http.Client
	url  string
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{},
		url:  "https://api.0x.org/swap/v1/",
	}
}

func (c *Client) sendSwapRequest(ctx context.Context, endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(c.url + endpoint)
	if err != nil {
		return "", err
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), strings.NewReader(""))
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(body), nil
}

func GetUSDSwapQuoteWithAmount(ctx context.Context, token, amount string) (*SwapQuote, error) {
	params := map[string]string{
		"sellAmount": amount,
		"buyToken":   "USDC",
		"sellToken":  token,
	}
	client := NewClient()
	body, err := client.sendSwapRequest(ctx, "quote", params)
	if err != nil {
		return nil, err
	}
	var quote SwapQuote
	err = json.Unmarshal([]byte(body), &quote)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func GetUSDSwapQuote(ctx context.Context, token string) (*SwapQuote, error) {
	const amount = "1000000000000000000" // Assume all tokens have 18 decimals
	params := map[string]string{
		"sellAmount": amount,
		"buyToken":   "USDC",
		"sellToken":  token,
	}
	client := NewClient()
	body, err := client.sendSwapRequest(ctx, "quote", params)
	if err != nil {
		return nil, err
	}
	var quote SwapQuote
	err = json.Unmarshal([]byte(body), &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

func GetETHSwapQuote(ctx context.Context, token string) (*SwapQuote, error) {
	const amount = "1000000000000000000"
	params := map[string]string{
		"sellAmount": amount,
		"buyToken":   "ETH",
		"sellToken":  token,
	}
	client := NewClient()
	body, err := client.sendSwapRequest(ctx, "quote", params)
	if err != nil {
		return nil, err
	}
	var quote SwapQuote
	err = json.Unmarshal([]byte(body), &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}
