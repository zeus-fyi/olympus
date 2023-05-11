package uniswap_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

type UniswapV2ApiToken struct {
	Id     string
	Symbol string
}

type UniswapV2ApiPair struct {
	Id       string
	Token0   UniswapV2ApiToken
	Token1   UniswapV2ApiToken
	Reserve0 *big.Int
	Reserve1 *big.Int
}

const Endpoint = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

func fetchPairs(ctx context.Context, graphqlQuery string) ([]UniswapV2ApiPair, error) {
	query := fmt.Sprintf(`{"query": %q}`, strings.ReplaceAll(graphqlQuery, "\n", " "))

	client := &http.Client{}
	req, err := http.NewRequest("POST", Endpoint, strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Pairs []struct {
				Id       string
				Token0   UniswapV2ApiToken
				Token1   UniswapV2ApiToken
				Reserve0 string
				Reserve1 string
			}
		}
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	pairs := make([]UniswapV2ApiPair, len(data.Data.Pairs))
	for i, pair := range data.Data.Pairs {
		reserve0 := new(big.Int)
		reserve0.SetString(pair.Reserve0, 10)

		reserve1 := new(big.Int)
		reserve1.SetString(pair.Reserve1, 10)

		pairs[i] = UniswapV2ApiPair{
			Id:       pair.Id,
			Token0:   pair.Token0,
			Token1:   pair.Token1,
			Reserve0: reserve0,
			Reserve1: reserve1,
		}
	}

	return pairs, nil
}

func GetTokenPairsWithVolume(ctx context.Context, limit, reserveMin, volumeMin int) ([]UniswapV2ApiPair, error) {
	graphqlQuery := fmt.Sprintf(`
        query {
            pairs(
                first: %d
                where: {
                    reserveUSD_gt: %d
                    volumeUSD_gt: %d
                }
                orderBy: reserveUSD
                orderDirection: asc
            ) {
                id
                token0 {
                    id
                    symbol
                }
                token1 {
                    id
                    symbol
                }
				reserve0
				reserve1
                reserveUSD
                volumeUSD

            }
        }
    `, limit, reserveMin, volumeMin)

	return fetchPairs(ctx, graphqlQuery)
}

func GetPairsForToken(ctx context.Context, limit int, tokenAddress string) ([]UniswapV2ApiPair, error) {
	graphqlQuery := fmt.Sprintf(`
		query {
			pairs(
				first: %d
				where: {
					or: [
						{ token0: "%s" }
						{ token1: "%s" }
					]
				}
				orderBy: reserveUSD
				orderDirection: asc
			) {
				id
				token0 {
					id
					symbol
				}
				token1 {
					id
					symbol
				}
				reserve0
				reserve1
				reserveUSD
				volumeUSD
			}
		}
	`, limit, tokenAddress, tokenAddress)

	return fetchPairs(ctx, graphqlQuery)
}
