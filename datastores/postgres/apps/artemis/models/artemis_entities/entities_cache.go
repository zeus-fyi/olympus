package artemis_entities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"golang.org/x/crypto/sha3"
)

type HashedRequestCache struct {
	RequestCache string `json:"requestCache"`
}

func HashWebRequestResultsAndParams(ou org_users.OrgUser, rt iris_models.RouteInfo) (*HashedRequestCache, error) {
	hp := []interface{}{rt.RoutePath, rt.RouteExt}
	if rt.Payloads != nil {
		hp = append(hp, rt.Payloads)
	} else if rt.Payload != nil {
		hp = append(hp, rt.Payload)
	}
	hash1, err := HashParams(ou.OrgID, hp)
	if err != nil {
		log.Err(err).Msg("Failed to hash request cache")
		return nil, err
	}
	return &HashedRequestCache{
		RequestCache: hash1,
	}, nil
}

func InsertEntitiesCaches(ctx context.Context, ue *UserEntityWrapper) (*HashedRequestCache, error) {
	err := InsertUserEntityLabeledMetadata(ctx, ue)
	if err != nil {
		log.Err(err).Msg("InsertEntitiesCaches: Failed to insert user entity")
		return nil, err
	}
	return nil, nil
}

func SelectEntitiesCaches(ctx context.Context, ue *UserEntityWrapper, ef EntitiesFilter) error {
	if ue == nil {
		return nil
	}
	err := SelectEntityLastMetadataData(ctx, ue, ef)
	if err != nil {
		return err
	}
	return nil
}

func SelectEntityLastMetadataData(ctx context.Context, ue *UserEntityWrapper, ef EntitiesFilter) error {
	if ue == nil {
		return nil
	}
	var qa string
	args := []interface{}{ue.Ou.OrgID, ef.Platform}
	if ef.SinceUnixTimestamp != 0 {
		args = append(args, ef.SinceUnixTimestamp)
		qa += fmt.Sprintf(" AND umd.entity_metadata_id >= $%d", len(args))
	}
	if ef.Nickname != "" {
		args = append(args, ef.Nickname)
		qa += fmt.Sprintf(" AND ue.nickname = $%d", len(args))
	}
	query := `
			SELECT umd.entity_id, umd.text_data, umd.json_data
            FROM public.user_entities ue
            JOIN public.user_entities_md umd ON ue.entity_id = umd.entity_id
            WHERE ue.org_id = $1 AND ue.platform = $2 ` + qa + `
			ORDER BY umd.entity_metadata_id DESC
			LIMIT 1;`

	um := UserEntityMetadata{}
	ue.MdSlice = []UserEntityMetadata{}
	um.JsonData = json.RawMessage{}
	err := apps.Pg.QueryRowWArgs(ctx, query, args...).Scan(&ue.EntityID, &um.TextData, &um.JsonData)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Err(err).Msg("SelectEntityLastMetadataDataWithLabel: Failed to execute query")
		return err
	}
	ue.MdSlice = append(ue.MdSlice, um)
	return nil
}

func HashParams(orgID int, hashParams []interface{}) (string, error) {
	hash := sha3.New256()
	for i, v := range hashParams {
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		if i == 0 {
			_, _ = hash.Write([]byte(fmt.Sprintf("org-%d", orgID)))
		}
		_, _ = hash.Write(b)
	}
	// Get the resulting encoded byte slice
	sha3v := hash.Sum(nil)
	return fmt.Sprintf("%x", hash.Sum(sha3v)), nil
}
