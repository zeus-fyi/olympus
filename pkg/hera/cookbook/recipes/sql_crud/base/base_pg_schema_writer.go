package base

import "strings"

func (m *ModelTemplate) WritePgTableDefinition() error {
	err := m.GetTableData()
	if err != nil {
		return err
	}

	for tbl, s := range m.StructMapToCodeGen {
		if strings.HasPrefix(tbl, "val") {
			continue
		}
		m.Path.AddGoFn(tbl)
		err = m.CreateTemplateFromStruct(s)
		if err != nil {
			return err
		}
	}
	return err
}
