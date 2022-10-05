package charts

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres"
)

func SelectPackageQuery(packageID int) string {
	query := fmt.Sprintf(`
	With cte_package as (
		SELECT cpk.chart_component_kind_name,
 			   cpk.chart_component_api_version,
  			   cp.chart_subcomponent_parent_class_type_id,
			   pct.chart_subcomponent_parent_class_type_name,
 			   cc.chart_subcomponent_child_class_type_name,
			   cc.chart_subcomponent_child_class_type_id,
			   cv.chart_subcomponent_key_name,
			   cv.chart_subcomponent_value,
			   cvj.chart_subcomponent_field_name,
			   cvj.chart_subcomponent_jsonb_key_values
		FROM chart_package_components cp
		LEFT JOIN chart_subcomponent_child_class_types cc ON cc.chart_subcomponent_parent_class_type_id = cp.chart_subcomponent_parent_class_type_id
		LEFT JOIN chart_subcomponents_child_values cv ON cc.chart_subcomponent_child_class_type_id = cv.chart_subcomponent_child_class_type_id
		LEFT JOIN chart_subcomponents_jsonb_child_values cvj ON cc.chart_subcomponent_child_class_type_id = cvj.chart_subcomponent_child_class_type_id
		INNER JOIN chart_subcomponent_parent_class_types pct ON pct.chart_subcomponent_parent_class_type_id = cp.chart_subcomponent_parent_class_type_id
		LEFT JOIN chart_component_kinds cpk ON cpk.chart_component_kind_id = pct.chart_component_kind_id
		WHERE cp.chart_package_id = %d
	)
	SELECT * FROM cte_package`, packageID)
	return query
}

type Package struct {
	chartComponentKindName   string
	chartComponentApiVersion string
	chartSubcomponents       PackageSubcomponent
}

type PackageSubcomponent struct {
	chartSubcomponentParentClassTypeId   int
	chartSubcomponentParentClassTypeName string
	chartSubcomponentChildClassTypeName  string
	chartSubcomponentChildClassTypeId    int
	chartSubcomponentKeyName             *string
	chartSubcomponentValue               *string
	chartSubcomponentFieldName           *string
	chartSubcomponentJsonbKeyValues      *string
}

type PackageComponentMap map[string]map[int][]PackageSubcomponent

func FetchQueryPackage(ctx context.Context, packageID int) (PackageComponentMap, error) {
	log.Info().Msg("FetchQueryPackage")

	query := SelectPackageQuery(packageID)
	packageComponents := make(PackageComponentMap)
	parentChildMap := make(map[int][]PackageSubcomponent)

	log.Debug().Interface("FetchQueryPackage: Query: ", query)
	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return packageComponents, err
	}
	defer rows.Close()
	for rows.Next() {
		var pkg Package
		var chartComp PackageSubcomponent
		rowErr := rows.Scan(
			&pkg.chartComponentKindName,
			&pkg.chartComponentApiVersion,
			&chartComp.chartSubcomponentParentClassTypeId,
			&chartComp.chartSubcomponentParentClassTypeName,
			&chartComp.chartSubcomponentChildClassTypeName,
			&chartComp.chartSubcomponentChildClassTypeId,
			&chartComp.chartSubcomponentKeyName,
			&chartComp.chartSubcomponentValue,
			&chartComp.chartSubcomponentFieldName,
			&chartComp.chartSubcomponentJsonbKeyValues,
		)
		if rowErr != nil {
			log.Err(rowErr).Interface("FindValidatorIndexes: Query: ", query)
			return packageComponents, rowErr
		}
		// maybe add apiVersion for more complex
		tmp := parentChildMap[chartComp.chartSubcomponentParentClassTypeId]
		tmp = append(tmp, chartComp)
		parentChildMap[chartComp.chartSubcomponentParentClassTypeId] = tmp
		packageComponents[pkg.chartComponentKindName] = parentChildMap
	}

	return packageComponents, err
}
