package artemis_entities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
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
	SinceTimeUnit      string          `json:"sinceTimeUnit,omitempty"`
	SinceTimeDuration  int             `json:"sinceTimeAmount,omitempty"`
}

func (ef *EntitiesFilter) SetSinceOffsetNowTimestamp(durationUnit string, durationAmount int) int {
	//ch := chronos.Chronos{}
	ts := time.Now() // Capturing the current time.
	// Adjusting ts based on the specified duration.
	switch durationUnit {
	case "minutes", "minute":
		ts = ts.Add(time.Duration(-durationAmount) * time.Minute)
	case "hours", "hour", "hrs", "hr":
		ts = ts.Add(time.Duration(-durationAmount) * time.Hour)
	case "days", "day":
		ts = ts.Add(time.Duration(-durationAmount) * 24 * time.Hour)
	case "weeks", "week":
		ts = ts.Add(time.Duration(-durationAmount) * 7 * 24 * time.Hour)
	default:
		// Handling unknown durationUnit appropriately.
	}
	ef.SinceUnixTimestamp = int(ts.UnixNano())
	return ef.SinceUnixTimestamp
}

func SelectUserMetadataByProvidedFields(ctx context.Context, ous org_users.OrgUser, nickname, platform string, labels []string, sinceUnixTimestamp int) ([]UserEntity, error) {
	var wrappers []UserEntity

	wc := "WHERE ue.org_id = $1"
	args := []interface{}{ous.OrgID}
	if len(nickname) > 0 {
		args = append(args, nickname)
		wc += fmt.Sprintf(" AND ue.nickname = $%d", len(args))
	}
	if len(platform) > 0 {
		args = append(args, platform)
		wc += fmt.Sprintf(" AND ue.platform = $%d", len(args))
	}
	if len(labels) > 0 {
		args = append(args, pq.Array(labels)) // Using pq.Array to ensure the slice is passed correctly
		wc += fmt.Sprintf(" AND umdl.label = ANY($%d)", len(args))
	}
	ch := chronos.Chronos{}
	if sinceUnixTimestamp > 0 || sinceUnixTimestamp < 0 {
		sinceUnixTimestamp = ch.AdjustedUnixTimestampNowRaw(sinceUnixTimestamp)
		args = append(args, sinceUnixTimestamp)
		wc += fmt.Sprintf(" AND umdl.entity_metadata_label_id > $%d", len(args))
	}
	baseQuery := `
		SELECT 
			ue.entity_id, 
			ue.nickname, 
			ue.platform, 
			ue.first_name, 
			ue.last_name, 
			umd.entity_metadata_id, 
			umd.json_data, 
			umd.text_data, 
			STRING_AGG(umdl.label, ', ') AS labels
        FROM public.user_entities ue
        LEFT JOIN public.user_entities_md umd ON ue.entity_id = umd.entity_id
        LEFT JOIN public.user_entities_md_labels umdl ON umd.entity_metadata_id = umdl.entity_metadata_id
        ` + wc + `
		GROUP BY 
			ue.entity_id, 
			ue.nickname, 
			ue.platform, 
			ue.first_name, 
			ue.last_name, 
			umd.entity_metadata_id, 
			umd.json_data, 
			umd.text_data `

	// Append ORDER BY clause at the end of the query
	finalQuery := baseQuery + " ORDER BY ue.platform DESC LIMIT 10000"

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
			userEntity           UserEntity
		)

		// Scan the row into local variables
		err = rows.Scan(&userEntity.EntityID, &userEntity.Nickname, &userEntity.Platform, &userEntity.FirstName, &userEntity.LastName, &metadataID, &jsonData, &textData, &label)
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

			// Split the label string and process if it's not nil
			if label != nil {
				labelsSlice := strings.Split(*label, ", ")
				// Process labelsSlice as needed
				for _, lbl := range labelsSlice {
					lbl = strings.TrimSpace(lbl) // Trim any leading/trailing space from each label
					metadataLabel := UserEntityMetadataLabel{
						Label: lbl,
					}
					metadata.Labels = append(metadata.Labels, metadataLabel)
				}
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
		wrappers = append(wrappers, wrapper.UserEntity)
	}

	return wrappers, nil
}
