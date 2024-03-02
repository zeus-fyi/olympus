package artemis_entities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type EntitiesFilter struct {
	Nickname           string          `json:"nickname" db:"nickname"`
	Platform           string          `json:"platform" db:"platform"`
	FirstName          *string         `json:"firstName,omitempty"`
	LastName           *string         `json:"lastName,omitempty"`
	Labels             []string        `json:"labels"`
	MetadataJsonb      json.RawMessage `json:"metadataJsonb,omitempty"`
	MetadataText       string          `json:"metadataText,omitempty"`
	SinceUnixTimestamp int             `json:"sinceTimestampUnix,omitempty"`
}

func SelectUserMetadataByNicknameAndPlatform(ctx context.Context, nickname, platform string, labels []string, sinceUnixTimestamp int) ([]UserEntityWrapper, error) {
	var wrappers []UserEntityWrapper

	baseQuery := `
        SELECT ue.entity_id, ue.nickname, ue.platform, ue.first_name, ue.last_name, umd.entity_metadata_id, umd.json_data, umd.text_data, umdl.label, umdl.entity_metadata_label_id
        FROM public.user_entities ue
        LEFT JOIN public.user_entities_md umd ON ue.entity_id = umd.entity_id
        LEFT JOIN public.user_entities_md_labels umdl ON umd.entity_metadata_id = umdl.entity_metadata_id
        WHERE ue.nickname = $1 AND ue.platform = $2`

	args := []interface{}{nickname, platform}
	argCounter := 2 // Keeps track of the argument index

	if len(labels) > 0 {
		argCounter++
		baseQuery += fmt.Sprintf(" AND umdl.label = ANY($%d)", argCounter)
		args = append(args, pq.Array(labels)) // Using pq.Array to ensure the slice is passed correctly
	}
	if sinceUnixTimestamp > 0 {
		argCounter++
		baseQuery += fmt.Sprintf(" AND umdl.entity_metadata_label_id > $%d", argCounter)
		args = append(args, sinceUnixTimestamp)
	}

	// Append ORDER BY clause at the end of the query
	finalQuery := baseQuery + " ORDER BY umdl.entity_metadata_label_id DESC"

	rows, err := apps.Pg.Query(ctx, finalQuery, args...)
	if err != nil {
		log.Err(err).Msg("SelectUserMetadataByNicknameAndPlatform: Failed to execute query")
		return nil, err
	}
	defer rows.Close()

	tempMap := make(map[int]*UserEntityWrapper)

	for rows.Next() {
		var (
			entityID, metadataID int
			jsonData             json.RawMessage
			textData, label      *string
			labelID              *int
			userEntity           UserEntity
		)

		// Scan the row into local variables
		err = rows.Scan(&entityID, &userEntity.Nickname, &userEntity.Platform, &userEntity.FirstName, &userEntity.LastName, &metadataID, &jsonData, &textData, &label, &labelID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan row")
			return nil, err
		}

		// Create or get the existing UserEntityWrapper from the map
		wrapper, exists := tempMap[entityID]
		if !exists {
			wrapper = &UserEntityWrapper{
				UserEntity: userEntity,
			}
			tempMap[entityID] = wrapper
		}

		// Append metadata and labels as necessary
		if metadataID != 0 {
			metadata := UserEntityMetadata{
				EntityMetadataID: metadataID,
				JsonData:         jsonData,
				TextData:         textData,
				Labels:           make([]UserEntityMetadataLabel, 0),
			}
			if label != nil && labelID != nil {
				metadataLabel := UserEntityMetadataLabel{
					EntityMetadataLabelID: *labelID,
					Label:                 *label,
				}
				metadata.Labels = append(metadata.Labels, metadataLabel)
			}
			wrapper.UserEntity.MdSlice = append(wrapper.UserEntity.MdSlice, metadata)
		}
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error occurred during rows iteration")
		return nil, err
	}

	// Convert map to slice
	for _, wrapper := range tempMap {
		wrappers = append(wrappers, *wrapper)
	}

	return wrappers, nil
}
