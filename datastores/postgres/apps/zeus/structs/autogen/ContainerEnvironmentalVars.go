package autogen_structs

type ContainerEnvironmentalVars struct {
	EnvID int    `db:"env_id"`
	Name  string `db:"name"`
	Value string `db:"value"`
}
