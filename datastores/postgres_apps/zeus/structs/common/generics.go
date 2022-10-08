package common

import autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"

type ChildValuesSlice []autogen_structs.ChartSubcomponentsChildValues
type ChildValuesMap map[string]autogen_structs.ChartSubcomponentsChildValues

func NewChildValuesMapKey(key string, m ChildValuesMap) ChildValuesMap {
	m[key] = autogen_structs.ChartSubcomponentsChildValues{}
	return m
}

func NewChildValuesMapKeyFromIterable(keys ...string) ChildValuesMap {
	m := ChildValuesMap{}
	for _, k := range keys {
		m = NewChildValuesMapKey(k, m)
	}
	return m
}
