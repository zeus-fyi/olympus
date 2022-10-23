package common

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type ChildValuesSlice []autogen_bases.ChartSubcomponentsChildValues
type ChildValuesMap map[string]autogen_bases.ChartSubcomponentsChildValues

func NewChildValuesMapKey(key string, m ChildValuesMap) ChildValuesMap {
	m[key] = autogen_bases.ChartSubcomponentsChildValues{}
	return m
}

func NewChildValuesMapKeyFromIterable(keys ...string) ChildValuesMap {
	m := ChildValuesMap{}
	for _, k := range keys {
		m = NewChildValuesMapKey(k, m)
	}
	return m
}
