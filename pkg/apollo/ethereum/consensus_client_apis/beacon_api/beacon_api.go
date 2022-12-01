package beacon_api

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	base_rest_client "github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var c client.Client

func init() {
	c.EnableBytesStrDecode = false
}

const getBlockByID = "eth/v2/beacon/blocks"
const getValidatorsByState = "eth/v1/beacon/states"

func GetValidatorsByState(ctx context.Context, beaconNode, stateID, status string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByState")
	if len(status) > 0 {
		status = fmt.Sprintf("?status=%s", status)
	}
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validators", status)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}

func GetValidatorsFinalized(ctx context.Context, beaconNode string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsFinalized")
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, "finalized/validators")
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}

func GetValidatorsByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL, status string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByStateFilter")
	if len(status) > 0 {
		status = fmt.Sprintf("&status=%s", status)
	}

	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validators?"+encodedQueryURL+status)
	log.Debug().Interface("BeaconAPI: url:", url)

	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}

func GetAllValidatorBalancesByState(ctx context.Context, beaconNode, stateID string) client.Reply {
	log.Info().Msg("BeaconAPI: GetAllValidatorBalancesByState")

	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validator_balances")
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}

func GetValidatorsBalancesByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsBalancesByStateFilter")

	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validator_balances?"+encodedQueryURL)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}

func GetBlockByID(ctx context.Context, beaconNode, blockID string) client.Reply {
	url := string_utils.UrlPathStrBuilder(getBlockByID, blockID)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	resp, err := r.R().Get(url)
	if err != nil {
		return client.Reply{}
	}
	reply := client.Reply{
		Body:       resp.String(),
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
		Err:        err,
		BodyBytes:  resp.Body(),
	}
	return reply
}
