package artemis_entities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type UserEntityWrapper struct {
	UserEntity
	Ou org_users.OrgUser `json:"-"`
}

type UserEntity struct {
	EntityID  int                  `json:"-" db:"entity_id"`
	Nickname  string               `json:"nickname" db:"nickname"`
	Platform  string               `json:"platform" db:"platform"`
	FirstName *string              `json:"firstName,omitempty" db:"first_name"` // Pointer used to handle NULL, omitempty for JSON if nil
	LastName  *string              `json:"lastName,omitempty" db:"last_name"`   // Pointer used to handle NULL, omitempty for JSON if nil
	MdSlice   []UserEntityMetadata `json:"metadata,omitempty" db:"metadata"`
}

type UserEntityMetadata struct {
	EntityMetadataID int                       `json:"-" db:"entity_metadata_id"`
	EntityID         int                       `json:"-" db:"entity_id"`
	JsonData         json.RawMessage           `json:"jsonData,omitempty" db:"json_data"` // Using json.RawMessage for JSONB
	TextData         *string                   `json:"textData,omitempty" db:"text_data"` // Pointer used to handle NULL, omitempty for JSON if nil
	Labels           []UserEntityMetadataLabel `json:"labels,omitempty" db:"labels"`
}

func (ue *UserEntity) GetStrLabels() []string {
	var lbs []string
	for _, mv := range ue.MdSlice {
		for _, lv := range mv.Labels {
			lbs = append(lbs, lv.Label)
		}
	}
	return lbs
}

type UserEntityMetadataLabel struct {
	EntityMetadataLabelID int    `json:"-" db:"entity_metadata_label_id"`
	EntityMetadataID      int    `json:"-" db:"entity_metadata_id"`
	Label                 string `json:"label" db:"label"`
}

func SearchLabelsForMatch(label string, ue UserEntity) bool {
	for _, mv := range ue.MdSlice {
		for _, lv := range mv.Labels {
			if lv.Label == label {
				return true
			}
		}
	}
	return false
}

func CreateMdLabels(labels []string) []UserEntityMetadataLabel {
	var lvs []UserEntityMetadataLabel
	for _, lv := range labels {
		lvs = append(lvs, UserEntityMetadataLabel{
			Label: lv,
		})
	}
	return lvs
}

func InsertUserEntityLabeledMetadata(ctx context.Context, ue *UserEntityWrapper) error {
	if ue == nil {
		return fmt.Errorf("nil UserEntityWrapper")
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Attempt to find or create the UserEntity based on nickname and platform
	findOrCreateUserQuery := `
	INSERT INTO public.user_entities (org_id, nickname, platform, first_name, last_name)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (org_id, nickname, platform) DO UPDATE
	SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
	RETURNING entity_id;`

	err = tx.QueryRow(ctx, findOrCreateUserQuery, ue.Ou.OrgID, ue.Nickname, ue.Platform, ue.FirstName, ue.LastName).Scan(&ue.EntityID)
	if err != nil {
		log.Printf("Failed to find or create user entity: %v", err)
		return err
	}

	// Loop through each UserEntityMetadata in MdSlice
	for mi, md := range ue.MdSlice {
		ue.MdSlice[mi].EntityID = ue.EntityID
		// Insert UserEntityMetadata using the obtained entityID
		insertMetadataQuery := `
		INSERT INTO public.user_entities_md (entity_id, json_data, text_data)
		VALUES ($1, $2, $3)
		RETURNING entity_metadata_id;`

		err = tx.QueryRow(ctx, insertMetadataQuery, ue.EntityID, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(md.JsonData), Status: IsNull(md.JsonData)}, md.TextData).Scan(&ue.MdSlice[mi].EntityMetadataID)
		if err != nil {
			log.Printf("Failed to insert user entity metadata: %v", err)
			return err
		}

		// Insert each UserEntityMetadataLabel for the current metadata
		insertLabelQuery := `
		INSERT INTO public.user_entities_md_labels (entity_metadata_label_id, entity_metadata_id, label)
		VALUES ($1, $2, $3)
		ON CONFLICT (entity_metadata_id, label) DO NOTHING
		RETURNING entity_metadata_label_id;`

		for i, label := range md.Labels {
			ch := chronos.Chronos{}
			err = tx.QueryRow(ctx, insertLabelQuery, ch.UnixTimeStampNow(), ue.MdSlice[mi].EntityMetadataID, label.Label).Scan(&md.Labels[i].EntityMetadataLabelID)
			if err != nil {
				log.Printf("Failed to insert user entity metadata label: %v", err)
				return err
			}
		}
	}
	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func IsNull(b []byte) pgtype.Status {
	if b == nil {
		return pgtype.Null
	}
	return pgtype.Present
}

func sanitizeBytesUTF8(b []byte) []byte {
	bs := bytes.ReplaceAll(b, []byte{0}, []byte{})
	return bs
}

func SelectHighestLabelIdForLabelAndPlatform(ctx context.Context, ou org_users.OrgUser, platform, label string) (int, error) {
	var highestLabelId *int

	query := `
		SELECT MAX(COALESCE(umdl.entity_metadata_label_id, 0)) 
		FROM public.user_entities_md_labels umdl
		JOIN public.user_entities_md umd ON umdl.entity_metadata_id = umd.entity_metadata_id
		JOIN public.user_entities ue ON umd.entity_id = ue.entity_id
		WHERE label = $1 AND platform = $2 AND ue.org_id = $3
		LIMIT 1;`

	err := apps.Pg.QueryRowWArgs(ctx, query, label, platform, ou.OrgID).Scan(&highestLabelId)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		log.Err(err).Msg("SelectHighestLabelIdForLabelAndPlatform: Failed to execute query")
		return 0, err
	}

	return aws.ToInt(highestLabelId), nil
}
