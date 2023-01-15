package validator_service_group

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
)

func InsertValidatorServiceOrgGroup(ctx context.Context, orgGroups hestia_autogen_bases.ValidatorServiceOrgGroupSlice, orgID int) error {
	tx, terr := apps.Pg.Begin(ctx)
	if terr != nil {
		log.Ctx(ctx).Err(terr)
		return fmt.Errorf("failed to start transaction: %v", terr)
	}

	for _, orgGroup := range orgGroups {
		_, err := tx.Exec(ctx, "INSERT INTO validators_service_org_groups (group_name, org_id, pubkey, protocol_network_id, fee_recipient, enabled) VALUES ($1, $2, $3, $4, $5, $6)", orgGroup.GroupName, orgID, orgGroup.Pubkey, orgGroup.ProtocolNetworkID, orgGroup.FeeRecipient, orgGroup.Enabled)
		if err != nil {
			log.Ctx(ctx).Err(err)
			rerr := tx.Rollback(ctx)
			if rerr != nil {
				log.Ctx(ctx).Err(rerr)
			}
			return fmt.Errorf("failed to insert into validators_service_org_groups: %v", err)
		}
	}
	err := tx.Commit(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}
