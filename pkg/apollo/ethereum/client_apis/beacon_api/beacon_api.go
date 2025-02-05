package beacon_api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func GetValidatorsByState(ctx context.Context, beaconNode, stateID, status string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("BeaconAPI: GetValidatorsByState")
	if len(status) > 0 {
		status = fmt.Sprintf("?status=%s", status)
	}
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validators", status)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	vsb := ValidatorsStateBeacon{}
	resp, err := r.R().SetResult(&vsb).Get(url)
	if err != nil {
		return vsb, err
	}
	if resp.StatusCode() != http.StatusOK {
		return vsb, errors.New("had a non-200 status code")
	}

	return vsb, err
}

func GetValidatorsFinalized(ctx context.Context, beaconNode string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("BeaconAPI: GetValidatorsFinalized")
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, "finalized/validators")
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	vsb := ValidatorsStateBeacon{}
	resp, err := r.R().SetResult(&vsb).Get(url)
	if err != nil {
		return vsb, err
	}
	if resp.StatusCode() != http.StatusOK {
		return vsb, errors.New("had a non-200 status code")
	}

	return vsb, err
}

func GetValidatorsByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL, status string) (ValidatorsStateBeacon, error) {
	log.Info().Msg("BeaconAPI: GetValidatorsByStateFilter")
	if len(status) > 0 {
		status = fmt.Sprintf("&status=%s", status)
	}

	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validators?"+encodedQueryURL+status)
	log.Debug().Interface("BeaconAPI: url:", url)

	bearer := "bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	r := base_rest_client.GetBaseRestyClient(beaconNode, bearer)
	vsb := ValidatorsStateBeacon{}
	resp, err := r.R().SetResult(&vsb).Get(url)
	if err != nil {
		return vsb, err
	}
	if resp.StatusCode() != http.StatusOK {
		return vsb, errors.New("had a non-200 status code")
	}

	return vsb, err
}

func GetAllValidatorBalancesByState(ctx context.Context, beaconNode, stateID string) (ValidatorBalances, error) {
	log.Info().Msg("BeaconAPI: GetAllValidatorBalancesByState")
	beaconNode = "https://CF62KTW23CWTUE2RNUFC:Q7PZEP72TGXYM2STKELBKMMAMDS6OSSG4GP42ZOT@mainnet.ethereum.coinbasecloud.net"
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validator_balances")
	log.Debug().Interface("BeaconAPI: url:", url)
	r := base_rest_client.GetBaseRestyClient(beaconNode, "")
	r.RetryCount = 3
	r.RetryWaitTime = 5 * 60 * time.Second

	vb := ValidatorBalances{}
	resp, err := r.R().SetResult(&vb).Get(url)
	if err != nil {
		return vb, err
	}
	if resp.StatusCode() != http.StatusOK {
		return vb, errors.New("had a non-200 status code")
	}

	return vb, err
}

func GetValidatorsBalancesByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) (ValidatorBalances, error) {
	log.Info().Msg("BeaconAPI: GetValidatorsBalancesByStateFilter")
	beaconNode = "https://CF62KTW23CWTUE2RNUFC:Q7PZEP72TGXYM2STKELBKMMAMDS6OSSG4GP42ZOT@mainnet.ethereum.coinbasecloud.net"
	url := string_utils.UrlPathStrBuilder(getValidatorsByState, stateID, "validator_balances?"+encodedQueryURL)
	log.Debug().Interface("BeaconAPI: url:", url)
	r := base_rest_client.GetBaseRestyClient(beaconNode, "")
	r.RetryCount = 10
	r.RetryWaitTime = 10 * time.Second
	vb := ValidatorBalances{}
	resp, err := r.R().SetResult(&vb).Get(url)
	if resp.StatusCode() != http.StatusOK {
		return vb, errors.New("had a non-200 status code")
	}
	return vb, err
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
