package autogen_structs

type ContainerProbes struct {
	ProbeID             int    `db:"probe_id"`
	ProbeKeyValuesJSONb string `db:"probe_key_values_jsonb"`
}
