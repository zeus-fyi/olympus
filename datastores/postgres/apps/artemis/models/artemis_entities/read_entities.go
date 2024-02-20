package artemis_entities

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type EntitiesFilter struct {
	Nickname           string          `json:"nickname" db:"nickname"`
	Platform           string          `json:"platform" db:"platform"`
	Labels             []string        `json:"labels"`
	MetadataJsonb      json.RawMessage `json:"metadataJsonb,omitempty"`
	MetadataText       string          `json:"metadataText,omitempty"`
	SinceUnixTimestamp int             `json:"sinceTimestampUnix,omitempty"`
}

func SelectUserMetadataByNicknameAndPlatform(ctx context.Context, nickname string, platform string, labels []string, sinceUnixTimestamp int) ([]UserEntityWrapper, error) {
	var wrappers []UserEntityWrapper

	// Construct the base query
	query := `
		SELECT ue.entity_id, ue.nickname, ue.platform, ue.first_name, ue.last_name, umd.entity_metadata_id, umd.json_data, umd.text_data, umdl.label, umdl.entity_metadata_label_id
		FROM public.user_entities ue
		LEFT JOIN public.user_entities_md umd ON ue.entity_id = umd.entity_id
		LEFT JOIN public.user_entities_md_labels umdl ON umd.entity_metadata_id = umdl.entity_metadata_id
		WHERE ue.nickname = $1 AND ue.platform = $2
		ORDER BY umdl.entity_metadata_label_id DESC`

	args := []interface{}{nickname, platform}

	// If labels are supplied, append them to the WHERE clause
	if len(labels) > 0 {
		query += " AND umdl.label = ANY($3)"
		args = append(args, labels)
	}
	if sinceUnixTimestamp > 0 {
		query += " AND umdl.entity_metadata_label_id > $4"
		args = append(args, sinceUnixTimestamp)
	}
	// Executing the query
	rows, err := apps.Pg.Query(ctx, query, args...)
	if err != nil {
		log.Err(err).Msg("SelectUserMetadataByNicknameAndPlatform: Failed to execute query")
		return nil, err
	}
	defer rows.Close()

	// Temporary map to hold unique UserEntityWrappers by entity_id
	tempMap := make(map[int]*UserEntityWrapper)

	for rows.Next() {
		var entityID, metadataID int
		var jsonData json.RawMessage
		var textData, label *string
		userEntity := UserEntity{
			MdSlice: []UserEntityMetadata{},
		}
		var metadata UserEntityMetadata
		var labelID *int
		err = rows.Scan(&entityID, &userEntity.Nickname, &userEntity.Platform, &userEntity.FirstName, &userEntity.LastName, &metadataID, &jsonData, &textData, &label, &labelID)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		// Check if we already have this UserEntity in our temp map
		wrapper, exists := tempMap[entityID]
		if !exists {
			wrapper = &UserEntityWrapper{
				UserEntity: userEntity,
			}
			tempMap[entityID] = wrapper
		}

		// Append metadata and labels as necessary
		if metadataID != 0 {
			metadata = UserEntityMetadata{
				EntityMetadataID: metadataID,
				JsonData:         jsonData,
				TextData:         textData,
				Labels:           []UserEntityMetadataLabel{},
			}
			if label != nil && labelID != nil {
				metadata.Labels = append(metadata.Labels, UserEntityMetadataLabel{
					EntityMetadataLabelID: *labelID,
					EntityMetadataID:      0,
					Label:                 *label,
				})
			}
			wrapper.MdSlice = append(wrapper.MdSlice, metadata)
		}
	}

	// Convert tempMap to slice
	for _, wrapper := range tempMap {
		wrappers = append(wrappers, *wrapper)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating retrieval rows: %v", err)
		return nil, err
	}

	return wrappers, nil
}
