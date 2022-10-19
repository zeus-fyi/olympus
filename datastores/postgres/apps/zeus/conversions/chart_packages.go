package conversions

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs"
)

func SelectPackageQuery(packageID int) string {
	query := fmt.Sprintf(`
	With cte_package as (
		SELECT cpr.chart_component_kind_name,
 			   cpr.chart_component_api_version,
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
		LEFT JOIN chart_component_resources cpr ON cpr.chart_component_resource_id = pct.chart_component_resource_id
		WHERE cp.chart_package_id = %d
	)
	SELECT * FROM cte_package`, packageID)
	return query
}

func FetchQueryPackage(ctx context.Context, packageID int) (structs.PackageComponentMap, error) {
	log.Info().Msg("FetchQueryPackage")

	query := SelectPackageQuery(packageID)
	packageComponents := make(structs.PackageComponentMap)
	parentChildMap := make(map[int][]structs.PackageSubcomponent)

	log.Debug().Interface("FetchQueryPackage: Query: ", query)
	rows, err := apps.Pg.Query(ctx, query)
	if err != nil {
		return packageComponents, err
	}
	defer rows.Close()
	for rows.Next() {
		var pkg structs.Package
		var chartComp structs.PackageSubcomponent
		rowErr := rows.Scan(
			&pkg.ChartComponentKindName,
			&pkg.ChartComponentApiVersion,
			&chartComp.ChartSubcomponentParentClassTypeId,
			&chartComp.ChartSubcomponentParentClassTypeName,
			&chartComp.ChartSubcomponentChildClassTypeName,
			&chartComp.ChartSubcomponentChildClassTypeId,
			&chartComp.ChartSubcomponentKeyName,
			&chartComp.ChartSubcomponentValue,
			&chartComp.ChartSubcomponentFieldName,
			&chartComp.ChartSubcomponentJsonbKeyValues,
		)
		if rowErr != nil {
			log.Err(rowErr).Interface("FindValidatorIndexes: Query: ", query)
			return packageComponents, rowErr
		}
		// maybe add apiVersion for more complex
		tmp := parentChildMap[chartComp.ChartSubcomponentParentClassTypeId]
		tmp = append(tmp, chartComp)
		parentChildMap[chartComp.ChartSubcomponentParentClassTypeId] = tmp
		packageComponents[pkg.ChartComponentKindName] = parentChildMap
	}

	return packageComponents, err
}

//cp := autogen_structs.ChartPackages{
//ChartPackageID: 0,
//ChartName:      "",
//ChartVersion:   "",
//}
//
//cpr := autogen_structs.ChartPackageComponents{
//ChartSubcomponentParentClassTypeID: 0,
//}
//cpk := autogen_structs.ChartComponentKinds{
//ChartComponentKindID:     0,
//ChartComponentKindName:   "",
//ChartComponentApiVersion: "",
//}
