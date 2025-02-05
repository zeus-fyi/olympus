package aegis_secrets

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertOrgSecretRef(ctx context.Context, orgSecretRef autogen_bases.OrgSecretReferences, secretRef autogen_bases.OrgSecretKeyValReferences) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH new_secret AS (
						INSERT INTO org_secret_references(org_id, secret_name, secret_id)
						VALUES($1, $2, $3) 
						ON CONFLICT DO NOTHING
					)
					INSERT INTO org_secret_key_val_references(secret_id, secret_env_var_ref, secret_key_ref, secret_name_ref)
					SELECT $3, $4, $5, $6 
				  	ON CONFLICT DO NOTHING;`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgSecretRef.OrgID, orgSecretRef.SecretName, orgSecretRef.SecretID, secretRef.SecretEnvVarRef, secretRef.SecretKeyRef, secretRef.SecretNameRef)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgSecretRef"))
}

func InsertOrgSecretTopologyRef(ctx context.Context, orgSecretTopRef autogen_bases.TopologySystemComponentsSecrets) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO topology_system_components_secrets(topology_system_component_id, secret_id)
				  SELECT $1, $2`
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgSecretTopRef.TopologySystemComponentID, orgSecretTopRef.SecretID)
	return misc.ReturnIfErr(err, q.LogHeader("InsertOrgSecretTopologyRef"))
}

type OrgSecretRef struct {
	autogen_bases.OrgSecretReferences
	autogen_bases.OrgSecretKeyValReferencesSlice
}

func DoesOrgSecretExistForTopology(ctx context.Context, orgID int, topName string) (bool, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT true
   				  FROM topology_system_components_secrets ts
   				  INNER JOIN topology_system_components tsys ON tsys.topology_system_component_id = ts.topology_system_component_id
   				  WHERE tsys.org_id = $1 AND tsys.topology_system_component_name = $2;`
	exists := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, topName).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return exists, misc.ReturnIfErr(err, q.LogHeader("DoesOrgSecretExistForTopology"))
}

func SelectOrgSecretRef(ctx context.Context, orgID int, topologyName string) (OrgSecretRef, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	
				  SELECT secret_name, secret_env_var_ref, secret_key_ref, secret_name_ref
				  FROM org_secret_references osr
				  INNER JOIN org_secret_key_val_references os ON os.secret_id = osr.secret_id
				  INNER JOIN topology_system_components_secrets ts ON ts.secret_id = osr.secret_id
				  INNER JOIN topology_system_components tsys ON tsys.topology_system_component_id = ts.topology_system_component_id
				  WHERE osr.org_id = $1 AND tsys.topology_system_component_name = $2;`
	log.Debug().Interface("SelectOrgSecretRef", q.LogHeader("SelectOrgSecretRef"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, topologyName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectOrgSecretRef")); returnErr != nil {
		return OrgSecretRef{}, err
	}
	orgRef := OrgSecretRef{}
	var orgSecretSlice autogen_bases.OrgSecretKeyValReferencesSlice
	defer rows.Close()
	for rows.Next() {
		secretRef := autogen_bases.OrgSecretKeyValReferences{}
		rowErr := rows.Scan(
			&orgRef.SecretName, &secretRef.SecretEnvVarRef, &secretRef.SecretKeyRef, &secretRef.SecretNameRef,
		)
		if rowErr != nil {
			return OrgSecretRef{}, rowErr
		}
		orgSecretSlice = append(orgSecretSlice, secretRef)
	}
	orgRef.OrgSecretKeyValReferencesSlice = orgSecretSlice
	return orgRef, misc.ReturnIfErr(err, q.LogHeader("SelectOrgSecretRef"))
}
