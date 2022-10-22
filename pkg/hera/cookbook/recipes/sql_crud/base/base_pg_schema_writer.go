package base

import "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives/structs"

func (m *ModelTemplate) WritePgTableDefinition() error {
	err := m.GetTableData()
	if err != nil {
		return err
	}

	for tbl, s := range m.StructMapToCodeGen {
		m.Structs = structs.NewStructsGen()
		m.Path.AddGoFn(tbl)
		err = m.CreateTemplateFromStruct(s)
		if err != nil {
			return err
		}
	}
	return err
}
