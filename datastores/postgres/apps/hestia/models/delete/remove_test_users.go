package hestia_delete

var q = `
WITH cte_clean_up_demo_org_users AS (
	SELECT ou.user_id as user_id, ou.org_id as org_id, quicknode_id AS quicknode_id
	FROM quicknode_marketplace_customer qmc 
	LEFT JOIN users_keys usk ON usk.public_key = qmc.quicknode_id
	LEFT JOIN org_users ou ON ou.user_id = usk.user_id
	WHERE is_test = true
	GROUP BY ou.user_id, ou.org_id, quicknode_id
), cte_users_keys AS (
	SELECT public_key
	FROM users_keys
	WHERE user_id IN (SELECT user_id FROM cte_clean_up_demo_org_users)
), cte_users_key_services AS ( 
	SELECT *
	FROM  users_key_services
	WHERE public_key IN (SELECT public_key FROM cte_users_keys)
)
  	SELECT pqs.quicknode_id as quicknode_id, pqs.endpoint_id as endpoint_id
	FROM provisioned_quicknode_services pqs
	INNER JOIN cte_clean_up_demo_org_users qc ON qc.quicknode_id = pqs.quicknode_id
`
