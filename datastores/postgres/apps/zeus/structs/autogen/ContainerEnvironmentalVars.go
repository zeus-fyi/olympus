package autogen_structs

type ContainerEnvironmentalVars struct {
	EnvID int                    `db:"env_id"`
	Name  string                 `db:"name"`
	Value map[string]interface{} `db:"value"`
}
