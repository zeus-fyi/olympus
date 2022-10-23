package base

import "github.com/zeus-fyi/olympus/pkg/utils/string_utils"

func (m *ModelTemplate) WritePgTableDefinition() error {
	err := m.GetTableData()
	if err != nil {
		return err
	}

	for tbl, s := range m.StructMapToCodeGen {
		if string_utils.FilterStringWithOpts(tbl, &m.Path.FilterFiles) {
			m.Path.AddGoFn(tbl)
			err = m.CreateTemplateFromStruct(s)
			if err != nil {
				return err
			}
		}
	}
	return err
}
