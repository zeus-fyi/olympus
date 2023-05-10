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

type UniswapToken struct {
	Id     string
	Symbol string
}

type UniswapPair struct {
	Id       string
	Token0   UniswapToken
	Token1   UniswapToken
	Reserve0 *big.Int
	Reserve1 *big.Int
}

const (
	Endpoint   = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
	SubgraphId = "QmWTrJJ9W8h3JE19FhCzzPYsJ2tgXZCdUqnbyuo64ToTBN"
)

func GetTokenPairsWithVolume(ctx context.Context) ([]UniswapPair, error) {
	graphqlQuery := `
        query {
            pairs(
                first: 1000
                where: {
                    reserveUSD_gt: 1000000
                    volumeUSD_gt: 50000
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
    `

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
				Token0   UniswapToken
				Token1   UniswapToken
				Reserve0 string
				Reserve1 string
			}
		}
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	pairs := make([]UniswapPair, len(data.Data.Pairs))
	for i, pair := range data.Data.Pairs {
		reserve0 := new(big.Int)
		reserve0.SetString(pair.Reserve0, 10)

		reserve1 := new(big.Int)
		reserve1.SetString(pair.Reserve1, 10)

		pairs[i] = UniswapPair{
			Id:       pair.Id,
			Token0:   pair.Token0,
			Token1:   pair.Token1,
			Reserve0: reserve0,
			Reserve1: reserve1,
		}
	}

	fmt.Println("Pairs 0:", data.Data.Pairs[0])
	return pairs, nil
}
