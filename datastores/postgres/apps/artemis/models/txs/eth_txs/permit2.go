package artemis_eth_txs

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (e *EthTx) GetPermit2Nonce() int {
	return 0
}

func (e *EthTx) PutPermit2Nonce() error {
	return nil
}

// TODO add permit2 nonce table

func insertPermit2Nonce(ctx context.Context, sd artemis_autogen_bases.EthScheduledDelivery) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO eth_tx(tx_hash, protocol_network_id)
				  VALUES ($1, $2)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, sd.PublicKey, sd.ProtocolNetworkID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("ArtemisPermit2Nonce")); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("ArtemisPermit2Nonce"))
}
