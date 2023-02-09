package hydra_eth2_web3signer

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"strconv"

	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

const (
	PHASE0    = "PHASE0"
	ALTAIR    = "ALTAIR"
	BELLATRIX = "BELLATRIX"
)

func DecodeBeaconBlockAndSlot(ctx context.Context, body any) (any, int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock")
		return nil, 0, err
	}
	m := make(map[string]any)
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
		return nil, 0, err
	}
	version := GetVersion(m)
	switch version {
	case PHASE0:
		blockPhase0 := consensys_eth2_openapi.BlockRequestPhase0{}
		err = json.Unmarshal(b, &blockPhase0)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
			return nil, 0, err
		}
		slot, serr := strconv.Atoi(blockPhase0.Block.Slot)
		if serr != nil {
			log.Ctx(ctx).Error().Err(serr).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
			return nil, 0, err
		}
		return blockPhase0, slot, err
	case ALTAIR:
		blockAltair := consensys_eth2_openapi.BlockRequestAltair{}
		err = json.Unmarshal(b, &blockAltair)
		if err != nil {
			return nil, 0, err
		}
		slot, serr := strconv.Atoi(blockAltair.Block.Slot)
		if serr != nil {
			log.Ctx(ctx).Error().Err(serr).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
			return nil, 0, err
		}
		return blockAltair, slot, err
	case BELLATRIX:
		blockBellatrix := consensys_eth2_openapi.BlockRequestBellatrix{}
		err = json.Unmarshal(b, &blockBellatrix)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
			return nil, 0, err
		}
		slot, serr := strconv.Atoi(blockBellatrix.BlockHeader.Slot)
		if serr != nil {
			log.Ctx(ctx).Error().Err(serr).Interface("body", body).Msg("DecodeBeaconBlock: Unmarshal")
			return nil, 0, err
		}
		return blockBellatrix, slot, err
	}
	return nil, 0, err
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
