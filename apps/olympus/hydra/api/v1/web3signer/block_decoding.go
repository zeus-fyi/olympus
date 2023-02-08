package hydra_eth2_web3signer

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"

	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

const (
	PHASE0    = "PHASE0"
	ALTAIR    = "ALTAIR"
	BELLATRIX = "BELLATRIX"
)

func DecodeBeaconBlock(ctx context.Context, body any) (any, error) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock")
		return nil, err
	}

	m := make(map[string]any)
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
		return nil, err
	}
	version := GetVersion(m)
	switch version {
	case PHASE0:
		blockPhase0 := consensys_eth2_openapi.BlockRequestPhase0{}
		err = json.Unmarshal(b, &blockPhase0)
		return blockPhase0, err
	case ALTAIR:
		blockAltair := consensys_eth2_openapi.BlockRequestAltair{}
		err = json.Unmarshal(b, &blockAltair)
		return blockAltair, err
	case BELLATRIX:
		blockBellatrix := consensys_eth2_openapi.BlockRequestBellatrix{}
		err = json.Unmarshal(b, &blockBellatrix)
		return blockBellatrix, err
	}
	return nil, err
}

func GetVersion(m map[string]any) string {
	for k, v := range m {
		if k == "version" {
			switch v {
			case PHASE0:
				return PHASE0
			case ALTAIR:
				return ALTAIR
			case BELLATRIX:
				return BELLATRIX
			default:
				// TODO
			}
		}
	}
	return ""
}
